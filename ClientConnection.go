package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type ClientConnection struct {
	myConn    net.Conn
	myEncoder *gob.Encoder
	myDecoder *gob.Decoder

	//tradeChannel is used to easily communicate with the client without going
	//directly through the event manager. When used, it is expected that the event
	//manager has set the characters trading flag.
	tradeChannel chan string

	//pingChannel is used to easily receive pings from the client.
	pingChannel chan string

	character *Character
}

//CliecntConnection constructor constructs a new client connection and starts
//a new thread to continuously read from the connection.
func NewClientConnection(conn net.Conn, char *Character, decoder *gob.Decoder, encder *gob.Encoder) {
	cc := new(ClientConnection)
	cc.myConn = conn

	cc.myEncoder = encder
	cc.myDecoder = decoder

	cc.character = char
	char.myClientConn = cc
	cc.character.myClientConn = cc

	cc.tradeChannel = make(chan string)
	cc.pingChannel = make(chan string)

	go cc.Read()
}

//Read will continuously try to read from the connection established with the client.
//If an error occurs or the client wishes to exit this will handle closing the connection
//and cleaning up the character from the world. When a succesful read occures
//the corresponding event is added to the event queu or executed rightaway.
func (cc *ClientConnection) Read() {
	defer cc.shutdown()

	for {
		var clientsMsg ClientMessage

		err := cc.myDecoder.Decode(&clientsMsg)
		checkError(err, false)

		fmt.Println("Message read: ", clientsMsg)

		if clientsMsg.CombatAction {
			eventManager.AddEvent(NewEvent(cc.character, clientsMsg.GetCommand(), clientsMsg.GetValue()))
		} else if clientsMsg.GetCommand() == "ping" {
			go cc.SendToPingChannel()
		} else if clientsMsg.IsTradeCommand() {
			go cc.SendToTradeChannel(clientsMsg)
		} else {
			go eventManager.ExecuteNonCombatEvent(cc, &clientsMsg)
		}

		if clientsMsg.Command == "exit" || err != nil {
			fmt.Println("Closing the connection")
			break
		}
	}
}

//shutdown ensures the clients player is removed from the world when the player
//disconnects, it also ensures the connection is closed.
func (cc *ClientConnection) shutdown() {
	eventManager.RemovePlayerFromRoom(cc.GetCharacter())
	cc.myConn.Close()
}

//Write will attempt to send the provided server message accross the connection
//to the client. Write automatically appends the clients most recent character
//data to the message.
func (cc *ClientConnection) Write(msg ServerMessage) {
	msg.addCharInfo(cc.character)
	err := cc.myEncoder.Encode(msg)
	checkError(err, false)
}

//SendToTradeChannel is used when the client is involved in a trade action with
//another player. When a clients wishes to add an item to the trade the item name
//is sent accross the trade channel so the event manager can add the items to the
//trade. To finish adding items the client must send the 'done' command
func (cc *ClientConnection) SendToTradeChannel(msg ClientMessage) {
	if msg.Command == "add" {
		for i := 0; i < msg.GetItemQuantity(); i++ {
			cc.tradeChannel <- msg.GetItem()
		}
	} else {
		cc.tradeChannel <- msg.GetValue()
	}
}

//GetResponseToTrade waits for the client to respond to some trade action, such
//as opening a trade or accepting trade terms. A timeout is used to prevent
//permantly blocking. Failure to respond within the time window is the same as
//rejecting the trade event.
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

func (cc *ClientConnection) GetItemsToTrade() string {
	return cc.GetResponseToTrade()
}

func (cc *ClientConnection) SendToPingChannel() {
	cc.pingChannel <- "ping"
}

func (cc *ClientConnection) GiveItem(itm Item_I) {
	cc.character.AddItem(itm)
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

func (cc *ClientConnection) GetAverageRoundTripTime() (avg time.Duration) {
	for i := 0; i < 10; i++ {
		cc.Write(newServerMessageTypeS(PING, "ping"))
		avg += cc.GetResponseToPing(time.Now())
	}

	return ((avg / 10) / 2)
}

func (cc *ClientConnection) GetCharactersName() string {
	return cc.character.Name
}

func (cc *ClientConnection) GetCharactersRoomID() int {
	return cc.character.RoomIN
}

func (cc *ClientConnection) GetCharacter() *Character {
	return cc.character
}

func (cc *ClientConnection) IsConnectionClosed() bool {
	return cc.myConn == nil
}
