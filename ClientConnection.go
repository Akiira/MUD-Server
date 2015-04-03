package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

type ClientConnection struct {
	myConn    net.Conn
	myEncoder *gob.Encoder
	myDecoder *gob.Decoder
	net_lock  sync.Mutex
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

	getCharactersFile(clientResponse.getUsername())
	cc.character = getCharacterFromFile(clientResponse.getUsername())
	cc.character.myClientConn = cc

	cc.CurrentEM = em

	//Send the client a description of their starting room
	em.executeNonCombatEvent(cc, &ClientMessage{Command: "look", Value: "room"})

	return cc
}

func (cc *ClientConnection) receiveMsgFromClient() {
	for {
		var clientResponse ClientMessage
		err := cc.myDecoder.Decode(&clientResponse)
		checkError(err, false)

		if err == nil {
			fmt.Println("Message read: ", clientResponse)
			if clientResponse.CombatAction {
				event := newEventFromMessage(clientResponse, cc.character, cc)
				cc.CurrentEM.addEvent(event)
			} else {
				cc.CurrentEM.executeNonCombatEvent(cc, &clientResponse)
			}

		} else {
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

func (cc *ClientConnection) getCharactersName() string {
	return cc.character.Name
}

func (cc *ClientConnection) getCharactersRoomID() int {
	return cc.character.RoomIN
}

func (cc *ClientConnection) getCharacter() *Character {
	return cc.character
}
