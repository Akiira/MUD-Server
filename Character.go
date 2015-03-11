package main

import (
	//"github.com/daviddengcn/go-colortext"
	"encoding/gob"
	"math/rand"
	"net"
	"sync"
)

// this should be a stub that hold a connection to a client
// works like a thread on its own
type Character struct {
	Name      string
	RoomIN    int
	HitPoints int
	Defense   int
	CurrentEM *EventManager
	myConn    net.Conn
	myEncoder *gob.Encoder
	myDecoder *gob.Decoder
	net_lock  sync.Mutex
	//messageQueue [100]ClientMessage

	//	Strength int
	//	Constitution int
	//	Dexterity int
	//	Wisdom int
	//	Charisma int
	//	Inteligence int

	//	Location string

	//	Race string
	//	Class string

	//	PersonalInvetory Inventory

	//	Weapon Item
	//ArmourSet
}

func (c *Character) init(conn net.Conn, name string, em *EventManager) {

	c.Name = name
	c.setCurrentEventManager(em)
	c.myEncoder = gob.NewEncoder(conn)
	c.myDecoder = gob.NewDecoder(conn)

}

func (c *Character) setCurrentEventManager(em *EventManager) {
	c.CurrentEM = em
}

func (c *Character) getEventMessage(msg ClientMessage) {
	//fmt.Print("I, ", (*c).Name, " receive msg : ")
	//fmt.Println(msg.Value)

	message := ClientMessage{Command: 1, Value: msg.Value}
	//c.net_lock.Lock()
	c.myEncoder.Encode(message)
	//c.net_lock.Unlock()

}

func (c *Character) receiveMessage() {

	var serversResponse ClientMessage
	for {
		c.net_lock.Lock()
		err := c.myDecoder.Decode(&serversResponse)
		c.net_lock.Unlock()
		checkError(err)
		//fmt.Println("message received")
		//fmt.Println(serversResponse.Value)
		if err == nil {
			c.CurrentEM.receiveMessage(serversResponse)
		}
	}
}

func (c *Character) getAttackRoll() int {
	return rand.Int() % 6
}

/*
func (c *Character) addItemToInventory(item Item) {

}

func (c *Character) equipItemFromGround(item Item) {

}

func (c *Character) equipItemFromInventory(itemName string) {

}
func (char *Character) moveCharacter(direction string) []FormattedString {
	room := worldRoomsG[char.RoomIN]
	dirAsInt := convertDirectionToInt(direction)

	if room.Exits[dirAsInt] >= 0 {
		room.removePCFromRoom(char.Name)
		room.ExitLinksToRooms[dirAsInt].addPCToRoom(char.Name)
		char.RoomIN = room.Exits[dirAsInt]
		return room.ExitLinksToRooms[dirAsInt].getFormattedOutput()
	} else {
		output := make([]FormattedString, 1, 1)
		output[0].Color = ct.Black
		output[0].Value = "No exit in that direction"
		return output
	}
}

func (c *Character) getAttack() int {
	return -1
}

func (c *Character) getName() string {
	return c.Name
}
*/
