package main

import (
	"encoding/gob"
	//"fmt"
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
}

//CliecntConnection constructor
func newClientConnection(conn net.Conn) *ClientConnection {
	cc := new(ClientConnection)
	cc.authen = true //need to be changed to false as default later
	cc.myConn = conn
	cc.myEncoder = gob.NewEncoder(conn)
	cc.myDecoder = gob.NewDecoder(conn)

	return cc
}

func (cc *ClientConnection) setConnection(conn net.Conn) {
	cc.myConn = conn
	cc.myEncoder = gob.NewEncoder(conn)
	cc.myDecoder = gob.NewDecoder(conn)
}

func (cc *ClientConnection) receiveMsgFromClient(em *EventManager) error {
	var clientResponse ClientMessage

	cc.net_lock.Lock()
	err := cc.myDecoder.Decode(&clientResponse)
	cc.net_lock.Unlock()
	//checkError(err)

	if err == nil {
		if(cc.authen){
			em.receiveMessage(clientResponse)	
		}
		else{
			//some logic to check for login or token
		}
		
	}

	return err

}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {

	//cc.net_lock.Lock()
	err := cc.myEncoder.Encode(msg)
	checkError(err)
	//cc.net_lock.Unlock()
}
