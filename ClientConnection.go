package main

import (
	"encoding/gob"
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
}

//CliecntConnection constructor
func newClientConnection(conn net.Conn) *ClientConnection {
	cc := new(ClientConnection)
	cc.myConn = conn
	cc.myEncoder = gob.NewEncoder(conn)
	cc.myDecoder = gob.NewDecoder(conn)

	//This associates the clients character with their connection
	var clientResponse ClientMessage
	err := cc.myDecoder.Decode(&clientResponse)
	checkError(err)

	cc.character = newCharacterFromName(clientResponse.Value)

	return cc
}

func (cc *ClientConnection) receiveMsgFromClient(em *EventManager) error {
	var clientResponse ClientMessage

	cc.net_lock.Lock()
	err := cc.myDecoder.Decode(&clientResponse)
	cc.net_lock.Unlock()
	//checkError(err)

	if err == nil {
		em.receiveMessage(clientResponse)
	}

	return err
}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {

	//cc.net_lock.Lock()
	err := cc.myEncoder.Encode(msg)
	checkError(err)
	//cc.net_lock.Unlock()
}
