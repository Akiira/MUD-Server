package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

type ClientConnection struct {
	myConn    net.Conn
	myEncoder *gob.Encoder
	myDecoder *gob.Decoder

	ping_lock    sync.Mutex
	pingResponse *sync.Cond
	tradeChannel chan string

	character *Character
	CurrentEM *EventManager
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

	cc.tradeChannel = make(chan string)
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
			cc.pingResponse.Signal()
		} else if clientResponse.IsTradeCommand() {
			cc.SendToTradeChannel(clientResponse)
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

func (cc *ClientConnection) GetItemsToTrade() string {
	return cc.GetResponseToTrade()
}

func (cc *ClientConnection) SendToTradeChannel(msg ClientMessage) {

	if msg.Command == "add" {
		for i := 0; i < msg.GetItemQuantity(); i++ {
			cc.tradeChannel <- msg.GetItem()
		}
	} else {
		cc.tradeChannel <- msg.GetValue()
	}
}

func (cc *ClientConnection) GetResponseToTrade() (response string) {
	timeoutChan := make(chan string)
	go func() {
		time.Sleep(time.Second * 60)
		timeoutChan <- "timeout"
	}()

	select {
	case msg := <-cc.tradeChannel:
		response = msg
	case msg := <-timeoutChan:
		response = msg
	}

	return response
}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {
	msg.addCharInfo(cc.character)
	err := cc.myEncoder.Encode(msg)
	checkError(err, false)
}

func (cc *ClientConnection) getAverageRoundTripTime() (avg time.Duration) {
	addr, _, err := net.SplitHostPort(cc.myConn.RemoteAddr().String())
	checkError(err, true)

	conn, err := net.Dial("tcp", addr+pingPort)
	checkErrorWithMessage(err, true, "Trying to open connection to client to get average round trip time.")
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for i := 0; i < 10; i++ {
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
