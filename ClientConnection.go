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
	pingChannel  chan string

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
	em.worldRooms[cc.character.RoomIN].AddPlayer(cc.character)
	cc.CurrentEM = em

	cc.pingResponse = sync.NewCond(&cc.ping_lock)

	//Send the client a description of their starting room
	em.executeNonCombatEvent(cc, &ClientMessage{Command: "look", Value: "room"})

	cc.tradeChannel = make(chan string)
	cc.pingChannel = make(chan string)
	return cc
}

func (cc *ClientConnection) receiveMsgFromClient() {
	defer cc.CurrentEM.RemovePlayerFromRoom(cc.getCharactersName(), cc.getCharactersRoomID())
	defer cc.myConn.Close()

	for {
		var clientResponse ClientMessage

		err := cc.myDecoder.Decode(&clientResponse)
		checkError(err, false)

		fmt.Println("Message read: ", clientResponse)

		if clientResponse.CombatAction {
			event := newEventFromMessage(clientResponse, cc.character)
			cc.CurrentEM.addEvent(event)
		} else if clientResponse.getCommand() == "ping" {
			cc.pingChannel <- "ping"
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
}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {
	msg.addCharInfo(cc.character)
	err := cc.myEncoder.Encode(msg)
	checkError(err, false)
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

func (cc *ClientConnection) GetResponseToPing(start time.Time) time.Duration {
	timeoutChan2 := make(chan string)
	go func() {
		time.Sleep(time.Second * 2)
		timeoutChan2 <- "timeout"
	}()

	select {
	case <-cc.pingChannel:
		return time.Now().Sub(start)
	case <-timeoutChan2:
		return time.Second * 6
	}
}

func (cc *ClientConnection) getAverageRoundTripTime() (avg time.Duration) {

	for i := 0; i < 10; i++ {
		cc.sendMsgToClient(newServerMessageTypeS(PING, "ping"))
		fmt.Println("\t\tWaiting for response ping")
		avg += cc.GetResponseToPing(time.Now())
	}

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
