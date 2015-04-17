package main

import (
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-colortext"
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

func (cc *ClientConnection) beginTrade(detail string) []FormattedString {
	var output []FormattedString
	if !cc.isTrading && !cc.myTrader.isSelected {
		fmt.Println(detail)
		dealerName := strings.Trim(detail, " ")
		dealer, found := cc.CurrentEM.worldRooms[cc.getCharactersRoomID()].CharactersInRoom[dealerName]

		if found && !dealer.getClientConnection().isTrading {
			fmt.Println("valid")
			cc.isTrading = true
			dealer.getClientConnection().isTrading = true
			cc.myTrader = &Trader{dealer: dealer.getClientConnection()}
			dealer.getClientConnection().dealerTrader = &Trader{dealer: cc}
			output = cc.character.PersonalInvetory.getInventoryDescription()
			str1 := "\nSelect item(s) to trade with command\n\"select itemIndex#1 quantity#1, itemIndex#2 quantity#2,..\""
			str2 := "\nEx. \"select 1 2, 4, 3 2\""
			str3 := "\nfor trade item#1 quantity:2, item#4 quantity:1, item#3 quantity:2"
			str4 := "\nYou can also reject with \"reject\" command"
			output = append(output, newFormattedString(str1+str2+str3+str4))

			invitation := dealer.PersonalInvetory.getInventoryDescription()
			str5 := "\nYou were invite to trade with " + cc.character.getName()
			invitation = append(invitation, newFormattedString(str5+str1+str2+str3+str4))
			dealer.getClientConnection().sendMsgToClient(newServerMessageTypeFS(GAMEPLAY, invitation))
			cc.myTrader.isSelected = true
			return output
		} else {
			fmt.Println("invalid")
			output = append(output, newFormattedString("\nCannot find "+dealerName+" in this room!\n"))
			return output
		}

	} else {
		fmt.Println("invalid")
		output = append(output, newFormattedString("\nYou are trading. You need to complete the current trade or reject it before start the new trading.\n"))
		return output
	}

}

func (cc *ClientConnection) selectItems(detail string) []FormattedString {
	var output []FormattedString
	if cc.isTrading {
		itemlist := strings.Split(detail, ",")
		var quan int
		var itemIndex int
		var validItem bool
		var name string
		var itemList []string
		var quanList []int
		for _, item := range itemlist {
			//fmt.Println(item)
			info := strings.Trim(item, " ")
			//fmt.Println(info)
			arguments := strings.Split(info, " ")
			if len(arguments) == 2 {
				//fmt.Println("[" + arguments[0] + "],[" + arguments[1] + "]")
				itemIndex, _ = strconv.Atoi(arguments[0])
				quan, _ = strconv.Atoi(arguments[1])
			} else {
				//fmt.Println("[" + arguments[0] + "]")
				itemIndex, _ = strconv.Atoi(arguments[0])
				quan = 1
			}

			validItem, name = cc.character.PersonalInvetory.checkAvailableItem(itemIndex, quan)
			if validItem {
				itemList = append(itemList, name)
				quanList = append(quanList, quan)
			} else {
				output = append(output, newFormattedString("\nYour item(s) selection is not valid. Please select again or reject\n"))
				return output
			}
		}

		output = append(output, newFormattedString("\nYou have selected the following item(s) to trade.\n"))
		var selectDeclare []FormattedString
		selectDeclare = append(selectDeclare, newFormattedString("\n"+cc.character.getName()+" have selected the following item(s) to trade.\n"))

		for i := 0; i < len(itemList); i++ {
			output = append(output, newFormattedString2(ct.Green, fmt.Sprintf("%d\t%-20s   %3d\n", (i+1), itemList[i], quanList[i])))
			selectDeclare = append(selectDeclare, newFormattedString2(ct.Green, fmt.Sprintf("%d\t%-20s   %3d\n", (i+1), itemList[i], quanList[i])))
		}

		cc.dealerTrader.dealer.sendMsgToClient(newServerMessageTypeFS(GAMEPLAY, selectDeclare))

		return output

	} else {
		fmt.Println("invalid")
		output = append(output, newFormattedString("\nYou are not trading. You need to initiate trading with command \"trade\" first.\n"))
		return output
	}
	return output

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
