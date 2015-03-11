package main

import (
	"encoding/gob"
	"github.com/daviddengcn/go-colortext"
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

	PersonalInvetory Inventory

	//	Weapon Item
	ArmourSet map[string]Armour
}

func newCharacter(name string, room int, hp int, def int) *Character {
	char := new(Character)
	char.Name = name
	char.HitPoints = hp
	char.Defense = def
	char.PersonalInvetory = *newInventory()
	char.ArmourSet = make(map[string]Armour)

	return char
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

//TODO change some of these functions so that they return []FormatterString
// 		so the client can see the effects.

func (c *Character) wearArmor(location string, armr Armour) {
	if _, ok := c.ArmourSet[location]; ok { // already an item present
		//TODO
	} else {
		c.ArmourSet[location] = armr
		c.Defense += armr.defense
	}
}

func (c *Character) takeOffArmor(location string) {
	if _, ok := c.ArmourSet[location]; ok { // already an item present
		delete(c.ArmourSet, location)
	} else {
		//TODO
	}
}

func (c *Character) addItemToInventory(item Item) {
	c.PersonalInvetory.items[item.name] = item
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
	if target.HP > 0 {
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
