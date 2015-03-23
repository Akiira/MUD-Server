package main

import (
	"encoding/xml"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"math/rand"
	"os"
)

// this should be a stub that hold a connection to a client
// works like a thread on its own
type Character struct {
	Name       string
	RoomIN     int
	HitPoints  int
	Defense    int
	level      int
	experience int

	Strength     int
	Constitution int
	Dexterity    int
	Wisdom       int
	Charisma     int
	Inteligence  int

	//	Race string
	//	Class string

	PersonalInvetory Inventory

	equipedWeapon  Weapon
	equippedArmour ArmourSet

	myClientConn *ClientConnection
}

func newCharacter(name string, room int, hp int, def int) *Character {
	char := new(Character)
	char.Name = name
	char.HitPoints = hp
	char.Defense = def
	char.PersonalInvetory = *newInventory()
	char.equippedArmour = newArmourSet()

	worldRoomsG[room].addPCToRoom(name)

	return char
}
func newCharacterFromName(name string) *Character {

	loadCharacterData(name)

	return onlinePlayers[name]
}

//TODO change some of these functions so that they return []FormatterString
// 		so the client can see the effects.

func (c *Character) wearArmor(location string, armr Armour) {
	if c.equippedArmour.isArmourEquippedAtLocation(location) { // already an item present
		//TODO
	} else {
		c.equippedArmour.equipArmour(location, armr)
	}
}

func (c *Character) takeOffArmor(location string) {
	if c.equippedArmour.isArmourEquippedAtLocation(location) { // already an item present
		c.equippedArmour.takeOffArmourByLocation(location)
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
		return room.ExitLinksToRooms[dirAsInt].getRoomDescription()
	} else {
		output := make([]FormattedString, 1, 1)
		output[0].Color = ct.Black
		output[0].Value = "No exit in that direction"
		return output
	}
}

func (char *Character) makeAttack(targetName string) []FormattedString {
	//TODO try to change this so it doesnt need global variable
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
	if target.HP <= 0 {
		// TODO  reward player exp
		output[1].Value = "\nThe " + targetName + " drops over dead."
		room := worldRoomsG[char.RoomIN]
		room.killOffMonster(targetName)
	}

	return output
}

func (c *Character) takeDamage(amount int, typeOfDamge int) []FormattedString {
	//TODO implement this
	return nil
}

func (c *Character) getAttackRoll() int {
	return rand.Int() % 6
}

func (c *Character) getDefense() int {
	return c.Defense
}

type CharacterXML struct {
	XMLName xml.Name `xml:"Character"`
	Name    string   `xml:"Name"`
	RoomIN  int      `xml:"RoomIN"`
	HP      int      `xml:"HitPoints"`
	Defense int      `xml:"Defense"`
}

func loadCharacterData(charName string) {
	//TODO remove hard coding
	xmlFile, err := os.Open("C:\\Go\\src\\MUD-Server\\Characters\\" + charName + ".xml")
	checkError(err)
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var charData CharacterXML
	xml.Unmarshal(XMLdata, &charData)

	char := newCharacter(charData.Name, charData.RoomIN, charData.HP, charData.Defense)
	onlinePlayers[charName] = char
}
