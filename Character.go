package main

import (
	"math/rand"
	"github.com/daviddengcn/go-colortext"
)

type Listener interface {
	//this will send a referebce to its own queue ti the eventmanager
	//so eventmanager can put broadcast msg into this object queue
	subscribeEventManager(EventManagerId int) bool
}

type Reporter interface {
	//this might probably just put msg onto a queue of eventmanager
	reportToEventManager(eventMsg ClientMessage) bool
}

// this should be a stub that hold a connection to a client
// works like a thread on its own
type Character struct {
	Name      string
	RoomIN    int
	HitPoints int
	Defense   int
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

func (c *Character) getAttackRoll() int {
	return rand.Int() % 6
}

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
