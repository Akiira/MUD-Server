package main

import (
	"encoding/gob"
	"fmt"
	//"fmt"
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
	fmt.Println("Connection: ", conn)
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
	fmt.Println("About to begin read loop in receiveMsgFromClient")
	for {
		cc.net_lock.Lock()
		fmt.Println("Aquired lock in receiveMsgFromClient")
		err := cc.myDecoder.Decode(&clientResponse)
		cc.net_lock.Unlock()
		fmt.Println("Released lock in receiveMsgFromClient")
		//checkError(err)

		if err == nil {
			fmt.Println("Received message from client: ", clientResponse)
			cc.CurrentEM.receiveMessage(clientResponse)
		} else {
			break
		}
	}

	fmt.Println("After read loop in receiveMsgFromClient")
}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {

	//cc.net_lock.Lock()
	err := cc.myEncoder.Encode(msg)
	checkError(err)
	//cc.net_lock.Unlock()
}

func (cc *ClientConnection) setCurrentEventManager(em *EventManager) {
	cc.CurrentEM = em

}
