package main

import (
	"encoding/gob"
	"fmt"
	"net"
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
}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {

	cc.net_lock.Lock()
	err := cc.myEncoder.Encode(msg)
	cc.net_lock.Unlock()
	checkError(err, false)
}

func (cc *ClientConnection) getAverageRoundTripTime() time.Duration {
	fmt.Println("\tGetting average round trip time.")

	var avg time.Duration
	for i := 0; i < 10; i++ {
		fmt.Println("\t\tPing: ", i)
		now := time.Now()
		cc.sendMsgToClient(newServerMessageTypeS(PING, "ping"))
		fmt.Println("\t\tWaiting for response ping")
		cc.pingResponse.L.Lock()
		cc.pingResponse.Wait()
		cc.pingResponse.L.Unlock()

		then := time.Now()
		avg += then.Sub(now)
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

func (cc *ClientConnection) giveItem(itm Item_I) {
	cc.character.addItemToInventory(itm)
}
