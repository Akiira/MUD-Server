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

	//send

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
	if !cc.isTrading {
		fmt.Println(detail)
		dealerName := strings.Trim(detail, " ")
		dealer, found := cc.CurrentEM.worldRooms[cc.getCharactersRoomID()].CharactersInRoom[dealerName]

		if found && !dealer.getClientConnection().isTrading {
			fmt.Println("valid")

			cc.isTrading = true
			dealer.getClientConnection().isTrading = true
			cc.myTrader = &Trader{isSelected: false, dealerCC: dealer.getClientConnection()}
			dealer.getClientConnection().myTrader = &Trader{isSelected: false, dealerCC: cc}

			output = cc.character.PersonalInvetory.getInventoryDescription()
			str1 := "\nSelect item(s) to trade with command\n\"select itemIndex#1 quantity#1, itemIndex#2 quantity#2,..\""
			str2 := "\nEx. \"select 1 2, 4, 3 2\""
			str3 := "\nfor trade item#1 quantity:2, item#4 quantity:1, item#3 quantity:2"
			str4 := "\nYou can also reject with \"reject\" command"
			output = append(output, newFormattedString(str1+str2+str3+str4))

			invitation := dealer.PersonalInvetory.getInventoryDescription()
			str5 := "\nYou were invited to trade with " + cc.character.getName()
			invitation = append(invitation, newFormattedString(str5+str1+str2+str3+str4))
			dealer.getClientConnection().sendMsgToClient(newServerMessageTypeFS(GAMEPLAY, invitation))

			return output
		} else {
			fmt.Println("invalid")
			output = append(output, newFormattedString("\nCannot find player \""+dealerName+"\" in this room!\n"))
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
	fmt.Println(cc.isTrading)
	fmt.Println(cc.myTrader)
	fmt.Println(cc.myTrader.isSelected)
	if cc.isTrading && cc.myTrader != nil && !cc.myTrader.isSelected {
		itemlist := strings.Split(detail, ",")
		var quan int
		var itemIndex int
		var validItem bool
		var name string
		itemMap := make(map[string]int)
		for _, item := range itemlist {

			//fmt.Println(item)
			info := strings.Trim(item, " ")
			if info == "" {
				continue
			}
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
				itemMap[name] = quan
			} else {
				output = append(output, newFormattedString("\nYour item(s) selection is not valid. Please select again or reject\n"))
				return output
			}
		}

		output = append(output, newFormattedString("\nYou have selected the following item(s) to trade.\n"))
		output = append(output, newFormattedString("\n---------------------------------------------------\n"))
		var selectDeclare []FormattedString
		selectDeclare = append(selectDeclare, newFormattedString("\n"+cc.character.getName()+" have selected the following item(s) to trade.\n"))
		selectDeclare = append(selectDeclare, newFormattedString("\n--------------------------------------------------------------------------\n"))
		i := 1
		for key, itemQuan := range itemMap {
			output = append(output, newFormattedString2(ct.Green, fmt.Sprintf("%d\t%-20s   %3d\n", i, key, itemQuan)))
			selectDeclare = append(selectDeclare, newFormattedString2(ct.Green, fmt.Sprintf("%d\t%-20s   %3d\n", i, key, itemQuan)))
			i++
		}

		output = append(output, newFormattedString("\n---------------------------------------------------\n"))
		selectDeclare = append(selectDeclare, newFormattedString("\n--------------------------------------------------------------------------\n"))

		selectDeclare = append(selectDeclare, newFormattedString2(ct.White, "You can \"accept\" or \"reject\"."))

		cc.myTrader.dealerCC.sendMsgToClient(newServerMessageTypeFS(GAMEPLAY, selectDeclare))

		cc.myTrader.itemMap = itemMap

		cc.myTrader.isSelected = true

		return output

	} else {
		fmt.Println("invalid")
		output = append(output, newFormattedString("\nYou are not trading. You need to initiate trading with command \"trade\" first.\n"))
		return output
	}
	return output

}

func (cc *ClientConnection) rejectTrading() []FormattedString {
	var output []FormattedString
	if cc.isTrading {

		output = append(output, newFormattedString("\nYou have canceled a trade with "+cc.myTrader.dealerCC.getCharactersName()+".\n"))

		var rejectDeclare []FormattedString
		rejectDeclare = append(rejectDeclare, newFormattedString2(ct.White, "\n"+cc.getCharactersName()+" have cenceled a trade with you.\n"))
		cc.myTrader.dealerCC.sendMsgToClient(newServerMessageTypeFS(GAMEPLAY, rejectDeclare))

		cc.myTrader.dealerCC.isTrading = false
		cc.myTrader.dealerCC.myTrader = nil

		cc.isTrading = false
		cc.myTrader = nil

		return output
	} else {
		fmt.Println("invalid")
		output = append(output, newFormattedString("\nYou are not trading.\n"))
		return output
	}
}

func (cc *ClientConnection) acceptTrading() []FormattedString {
	var output []FormattedString
	if cc.isTrading && cc.myTrader.isSelected {

		if cc.myTrader.dealerCC.myTrader.isConfirmed {
			//TODO finalize trading
			result := Trading(cc.myTrader, cc.myTrader.dealerCC.myTrader)

			if result {
				output = append(output, newFormattedString("Your trade with "+cc.myTrader.dealerCC.getCharactersName()+" has been confirmed.\n"))
			} else {
				output = append(output, newFormattedString("Your trade with "+cc.myTrader.dealerCC.getCharactersName()+" has been canceled.\n"))
			}

			ccXML := cc.character.toXML()
			ccXML.CurrentWorld = serverName
			sendCharactersXML(ccXML)

			ccXML2 := cc.myTrader.dealerCC.character.toXML()
			ccXML2.CurrentWorld = serverName
			sendCharactersXML(ccXML2)

			cc.myTrader.dealerCC.isTrading = false
			cc.myTrader.dealerCC.myTrader = nil
			cc.isTrading = false
			cc.myTrader = nil

			return output
		} else {
			cc.myTrader.isConfirmed = true
			output = append(output, newFormattedString("Please wait for "+cc.myTrader.dealerCC.getCharactersName()+" to confirm this trade.\n"))

			var acceptDeclare []FormattedString
			acceptDeclare = append(acceptDeclare, newFormattedString2(ct.White, "\n"+cc.getCharactersName()+" have confirmed this trade with you.\nPlease confirm or reject this trade.\n"))
			cc.myTrader.dealerCC.sendMsgToClient(newServerMessageTypeFS(GAMEPLAY, acceptDeclare))

			return output
		}
	} else {
		fmt.Println("invalid")
		output = append(output, newFormattedString("\nYou cannot accept yet. You need to select zero or more items to trade first.\n"))
		return output
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
