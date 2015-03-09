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

	PersonalInvetory Inventory

	//	Weapon Item
	//ArmourSet
}

func newCharacter(name string, room int, hp int, def int) *Character {
	char := new(Character)
	char.Name = name
	char.HitPoints = hp
	char.Defense = def
	char.PersonalInvetory = *newInventory()
	
	return char
}

//func newCharacterFromCharacter(oldChar Character) *Character {
//	char := new(Character)
//	char.Name = oldChar.Name
//	char.HitPoints = oldChar.HitPoints
//	char.Defense = oldChar.Defense
//}

func (c *Character) getAttackRoll() int {
	return rand.Int() % 6
}

func (c *Character) addItemToInventory(item Item) {
	c.PersonalInvetory.items[item.name] = item
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

func (char *Character) makeAttack(targetName string) []FormattedString {
	target := worldRoomsG[char.RoomIN].getMonster(targetName)
	output := make([]FormattedString, 2, 2)

	a1 := char.getAttackRoll()
	if a1 >= target.Defense {
		target.HP -= 2
		output[0].Value = "\nYou hit the " + targetName + "!"
	} else {
		output[0].Value = "\nYou missed the " + targetName + "!"
	}

	a2 := target.getAttackRoll()
	if ( target.HP > 0 ) {
		if a2 >= char.Defense {
			char.HitPoints -= 1
			output[1].Value = "\nThe " + targetName + " hit you!"
		} else {
			output[1].Value = "\nThe " + targetName + " narrowly misses you!"
		}	
	} else { //TODO add corpse to Rooms list of items
			// TODO  reward player exp
		output[1].Value = "\nThe " + targetName + " drops over dead."
		room := worldRoomsG[char.RoomIN]
		room.killOffMonster(targetName)
	}
	
	return output
}



func (c *Character) getName() string {
	return c.Name
}
