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
		em.receiveMessage(clientResponse)
	}

	return err

}

func (cc *ClientConnection) sendMsgToClient(msg ServerMessage) {
	//fmt.Println(cc.myEncoder)
	//cc.net_lock.Lock()
	err := cc.myEncoder.Encode(msg)
	if err != nil {
		//fmt.Printf("I detect error at send")
		//fmt.Println(err)
	}
	//cc.net_lock.Unlock()
}
