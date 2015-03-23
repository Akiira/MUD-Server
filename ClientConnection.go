package main

import (
	"encoding/gob"
	_ "fmt"
	//"io"
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
	checkError(err) //TODO replace check errors with somthing that doesnt crash server

	cc.character = newCharacterFromName(clientResponse.Value)
	cc.CurrentEM = em

	startingRoomDescription := worldRoomsG[cc.character.RoomIN].getRoomDescription()
	err = cc.myEncoder.Encode(ServerMessage{Value: startingRoomDescription})
	checkError(err)

	return cc
}

func (cc *ClientConnection) receiveMsgFromClient() {
	var clientResponse ClientMessage

	for {
		err := cc.myDecoder.Decode(&clientResponse)
		checkError(err)

		if err == nil {
			if clientResponse.combatAction {
				event = newEvent(PLAYER, cc.character.Name, clientResponse.Command, clientResponse.Value)
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
	checkError(err)
}

func (cc *ClientConnection) setCurrentEventManager(em *EventManager) {
	cc.CurrentEM = em
}
