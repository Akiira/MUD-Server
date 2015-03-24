package main

import (
	"encoding/gob"
	"fmt"
	//"io"
	"net"
	"sync"
)

type ClientConnection struct {
	authen    bool
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
	cc.authen = true //need to be changed to false as default later
	cc.myConn = conn

	cc.myEncoder = gob.NewEncoder(conn)
	cc.myDecoder = gob.NewDecoder(conn)

	//This associates the clients character with their connection
	var clientResponse ClientMessage
	err := cc.myDecoder.Decode(&clientResponse)
	checkError(err) //TODO replace check errors with somthing that doesnt crash server

	if err == nil {
		if(cc.authen){
			em.receiveMessage(clientResponse)	
		}
		else{
			//some logic to check for login or token
		}
		
	}
	
	cc.character = newCharacterFromName(clientResponse.Value)
	cc.CurrentEM = em

	startingRoomDescription := worldRoomsG[cc.character.RoomIN].getRoomDescription()
	err = cc.myEncoder.Encode(ServerMessage{Value: startingRoomDescription})
	checkError(err)

	return cc
}

func (cc *ClientConnection) receiveMsgFromClient() {
	for {
		var clientResponse ClientMessage
		err := cc.myDecoder.Decode(&clientResponse)
		checkError(err)

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
	checkError(err)
}

func (cc *ClientConnection) getCharactersName() string {
	return cc.character.Name
}
func (cc *ClientConnection) setCurrentEventManager(em *EventManager) {
	cc.CurrentEM = em
}
