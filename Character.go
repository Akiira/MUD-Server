package main

import (
	"encoding/xml"
	//"fmt"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"math/rand"
	"os"
)

type Character struct {
	Name   string
	RoomIN int

	HitPoints    int
	MaxHitPoints int
	Defense      int
	level        int
	experience   int

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

	return char
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

//TODO
//func (char *Character) moveCharacter(direction string, source *Room, destination *Room) []FormattedString

func (char *Character) moveCharacter(direction string) []FormattedString {
	//TODO this is just a temporary fix
	room := char.myClientConn.CurrentEM.worldRooms[char.RoomIN]

	dirAsInt := convertDirectionToInt(direction)

	if room.Exits[dirAsInt] >= 0 {
		room.removePCFromRoom(char.Name)
		room.ExitLinksToRooms[dirAsInt].addPCToRoom(char)
		char.RoomIN = room.Exits[dirAsInt]
		return room.ExitLinksToRooms[dirAsInt].getRoomDescription()
	} else {
		output := make([]FormattedString, 1, 1)
		output[0].Color = ct.Black
		output[0].Value = "No exit in that direction"
		return output
	}
}

//TODO
//func (char *Character) makeAttack(target *Agenter) []FormattedString

func (char *Character) makeAttack(targetName string) []FormattedString {
	//TODO again just a temporary fix
	target := char.myClientConn.CurrentEM.worldRooms[char.RoomIN].getMonster(targetName)

	output := make([]FormattedString, 2, 2)

	a1 := char.getAttackRoll()
	if a1 >= target.Defense {
		target.takeDamage(2, 0)
		output[0].Value = "\nYou hit the " + targetName + "!"
	} else {
		output[0].Value = "\nYou missed the " + targetName + "!"
	}

	if target.HP <= 0 {
		// TODO  reward player exp
		output[1].Value = "\nThe " + targetName + " drops over dead."
		room := char.myClientConn.CurrentEM.worldRooms[char.RoomIN]
		room.killOffMonster(targetName)
	}

	return output
}

func (c *Character) takeDamage(amount int, typeOfDamge int) []FormattedString {
	//TODO implement this
	return nil
}

func (c *Character) getAttackRoll() int {
	return rand.Int() % 20
}

func (c *Character) getDefense() int {
	return c.Defense
}

//func (c *Character) getStatsPage() []FormattedString {
//	stats := make([]FormattedString, 10, 20)

//	stats[0].Color = ct.Green
//	stats[0].Value = "Character Page for " + c.Name + "-------------------------------------------------\n"

//	stats[1].Value = "LEVEL:"
//	stats[2].Value = fmt.Sprintf("%2d %8s", c.level, "")
//	stats[3].Value = "RACE :"
//	stats[4].Value = fmt.Sprintf("%8s\n", "Human") //TODO
//	stats[5].Value = "AGE  :"
//	stats[6].Value = fmt.Sprintf("%4d %6s", 123, "") //TODO
//	stats[7].Value = "CLASS:"
//	stats[8].Value = fmt.Sprintf("%8s\n", "Ranger") //TODO
//	stats[9].Value = "STR  :"
//	stats[10].Value = fmt.Sprintf("%2d %8s", c.Strength, "")
//	stats[11].Value = "HitRoll:"
//	stats[12].Value = fmt.Sprintf("%8s\n", "66") //TODO
//	stats[13].Value = "INT  :"
//	stats[14].Value = fmt.Sprintf("%2d %8s", c.Inteligence, "")
//	stats[15].Value = "DmgRoll:"
//	stats[16].Value = fmt.Sprintf("%8s\n", "66") //TODO
//	stats[17].Value = "WIS  :"
//	stats[18].Value = fmt.Sprintf("%2d %8s", c.Strength, "")
//	stats[19].Value = "Alignment:"
//	stats[20].Value = fmt.Sprintf("%8s\n", "Paragon") //TODO
//	stats[21].Value = "DEX  :"
//	stats[22].Value = fmt.Sprintf("%2d %8s", c.Wisdom, "")
//	stats[23].Value = "Armour:"
//	stats[24].Value = fmt.Sprintf("%8s\n", "-500") //TODO
//	stats[25].Value = "CON  :"
//	stats[26].Value = fmt.Sprintf("%2d %8s", c.Constitution, "")

//	stats[1].Value = "CHA  :"
//	stats[2].Value = fmt.Sprintf("%2d %8s", c.Charisma, "")
//}

type CharacterXML struct {
	XMLName      xml.Name `xml:"Character"`
	Name         string   `xml:"Name"`
	RoomIN       int      `xml:"RoomIN"`
	HP           int      `xml:"HitPoints"`
	Defense      int      `xml:"Defense"`
	Password     string   `xml:"Password"`
	CurrentWorld string   `xml:"CurrentWorld"`
}

func getCharacterFromFile(charName string) *Character {
	//TODO remove hard coding
	xmlFile, err := os.Open("C:\\Go\\src\\MUD-Server\\Characters\\" + charName + ".xml")
	checkError(err)
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var charData CharacterXML
	xml.Unmarshal(XMLdata, &charData)

	char := newCharacter(charData.Name, charData.RoomIN, charData.HP, charData.Defense)

	return char
}
