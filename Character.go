package main

import (
	"io"
	//"github.com/daviddengcn/go-colortext"
	"math/rand"
	"net"
)

// this should be a stub that hold a connection to a client
// works like a thread on its own
type Character struct {
	Name         string
	RoomIN       int
	HitPoints    int
	Defense      int
	CurrentEM    *EventManager
	myClientConn *ClientConnection

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
	c.myClientConn = new(ClientConnection)
	c.myClientConn.setConnection(conn)

}

func (c *Character) setCurrentEventManager(em *EventManager) {
	c.CurrentEM = em

}

func (c *Character) getEventMessage(msg ServerMessage) {
	//fmt.Print("I, ", (*c).Name, " receive msg : ")
	//fmt.Println(msg.Value)
	c.myClientConn.sendMsgToClient(msg)

}

func (c *Character) receiveMessage() {

	go c.routineReceiveMsg()
}

func (c *Character) routineReceiveMsg() {

	for {
		err := c.myClientConn.receiveMsgFromClient(c.CurrentEM)
		if err == io.EOF {
			//need to unsubscribe and let this character be devour by garbage collecter
			c.CurrentEM.unsubscribeListener(c)
			break
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
