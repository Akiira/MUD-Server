package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ClientConnection struct {
	myConn       net.Conn
	myEncoder    *gob.Encoder
	myDecoder    *gob.Decoder
	net_lock     sync.Mutex
	ping_lock    sync.Mutex
	pingResponse *sync.Cond
	character    *Character
	CurrentEM    *EventManager
	isTrading    bool
	isOfferer    bool
	myTrader     *Trader
	dealerTrader *Trader
}

//CliecntConnection constructor
func newClientConnection(conn net.Conn, em *EventManager) *ClientConnection {
	cc := new(ClientConnection)
	cc.myConn = conn

	cc.myEncoder = gob.NewEncoder(conn)
	cc.myDecoder = gob.NewDecoder(conn)

	//This associates the clients character with their connection
	var clientResponse ClientMessage
	err := cc.myDecoder.Decode(&clientResponse)
	checkError(err, true)

	cc.character = getCharacterFromCentral(clientResponse.getUsername())
	cc.character.myClientConn = cc
	em.worldRooms[cc.character.RoomIN].addPCToRoom(cc.character)
	cc.CurrentEM = em

	cc.pingResponse = sync.NewCond(&cc.ping_lock)

	//Send the client a description of their starting room
	em.executeNonCombatEvent(cc, &ClientMessage{Command: "look", Value: "room"})

	return cc
}

func (cc *ClientConnection) receiveMsgFromClient() {

	for {
		var clientResponse ClientMessage

		err := cc.myDecoder.Decode(&clientResponse)
		checkError(err, false)

		fmt.Println("Message read: ", clientResponse)

		if clientResponse.CombatAction {
			event := newEventFromMessage(clientResponse, cc.character)
			cc.CurrentEM.addEvent(event)
		} else if clientResponse.getCommand() == "ping" {
			fmt.Println("\t\tReceived ping from user.")
			cc.pingResponse.Signal()
		} else {
			cc.CurrentEM.executeNonCombatEvent(cc, &clientResponse)
		}

		if clientResponse.Command == "exit" || err != nil {
			fmt.Println("Closing the connection")
			break
		}
	}

	cc.myConn.Close()
	cc.myConn = nil
}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {

	cc.net_lock.Lock()
	msg.addCharInfo(cc.character)
	err := cc.myEncoder.Encode(msg)
	cc.net_lock.Unlock()
	checkError(err, false)
}

func (cc *ClientConnection) getAverageRoundTripTime() time.Duration {
	fmt.Println("\tGetting average round trip time.")

	var avg time.Duration
	addr, _, err := net.SplitHostPort(cc.myConn.RemoteAddr().String())
	checkError(err, true)
	conn, err := net.Dial("tcp", addr+pingPort)
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for i := 0; i < 10; i++ {
		fmt.Println("\t\tPing: ", i)
		now := time.Now()
		err = encoder.Encode(newServerMessageS("ping"))
		checkError(err, false)
		if err != nil {
			avg += time.Minute * 10
			break
		}
		fmt.Println("\t\tWaiting for response ping")
		err = decoder.Decode(newClientMessage("", ""))
		checkError(err, false)
		then := time.Now()

		if err != nil {
			avg += time.Minute * 10
			break
		}
		fmt.Println("Time diff: ", then.Sub(now))
		avg += then.Sub(now)
	}
	encoder.Encode(newServerMessageS("done"))
	fmt.Println("\tDone getting average round trip time.")
	return ((avg / 10) / 2)
}

func (cc *ClientConnection) beginTrade(detail string) {
	if !cc.isTrading {
		fmt.Println(detail)
		if cc.character.PersonalInvetory.isValidTradeCmd(detail) {

			arguments := strings.Split(detail, " ")
			dealerName := arguments[1]
			var itemNum int
			dealer, found := cc.CurrentEM.worldRooms[cc.getCharactersRoomID()].CharactersInRoom[dealerName]
			if len(arguments) == 3 {
				itemNum, _ = strconv.Atoi(arguments[2])
			} else {
				itemNum = 1
			}
			if found && !dealer.getClientConnection().isTrading {
				fmt.Println("valid")
				cc.isTrading = true
				dealer.getClientConnection().isTrading = true
				cc.myTrader = &Trader{itemName: arguments[0], quantity: itemNum, dealer: dealer.getClientConnection()}
				dealer.getClientConnection().dealerTrader = &Trader{itemName: arguments[0], quantity: itemNum, dealer: cc}
			} else {
				fmt.Println("invalid")
			}

		} else {
			fmt.Println("invalid")
		}

	}
}

func (cc *ClientConnection) getCharactersName() string {
	return cc.character.Name
}

func (cc *ClientConnection) getCharactersRoomID() int {
	return cc.character.RoomIN
}

func (cc *ClientConnection) getCharacter() *Character {
	return cc.character
}

func (cc *ClientConnection) isConnectionClosed() bool {
	return cc.myConn == nil
}

func (cc *ClientConnection) giveItem(itm Item_I) {
	cc.character.addItemToInventory(itm)
}
