package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"math/rand"
	"os"
)

type Character struct {
	Race  string
	Class string

	Defense    int
	level      int
	experience int

	Agent

	//Equipment fields
	PersonalInvetory Inventory
	equipedWeapon    Weapon
	equippedArmour   ArmourSet

	myClientConn *ClientConnection
}

func newCharacter(name string, room int, hp int, def int) *Character {
	char := new(Character)
	char.Name = name
	char.currentHP = hp
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

func (char *Character) makeAttack(target Agenter) []FormattedString {

	output := make([]FormattedString, 2, 2)

	a1 := char.getAttackRoll()
	if a1 >= target.getDefense() {
		target.takeDamage(2, 0)
		output[0].Value = "\nYou hit the " + target.getName() + "!"
	} else {
		output[0].Value = "\nYou missed the " + target.getName() + "!"
	}

	if target.isDead() {
		// TODO  reward player exp
		output[1].Value = "\nThe " + target.getName() + " drops over dead."
		room := char.myClientConn.CurrentEM.worldRooms[char.RoomIN] //TODO fix this line
		room.killOffMonster(target.getName())
	}

	return output
}

func (c *Character) takeDamage(amount int, typeOfDamge int) []FormattedString {
	if amount-c.equippedArmour.getArmoursDefense() > 0 {
		c.currentHP -= (amount - c.equippedArmour.getArmoursDefense())
	}
	s := "You got hit for " + fmt.Sprintf("%i", amount) + " damage.\n"
	return newFormattedStringSplice2(ct.Red, s)
}
func (c *Character) isDead() bool {
	return c.currentHP > 0
}
func (c *Character) getName() string {
	return c.Name
}
func (c *Character) getRoomID() int {
	return c.RoomIN
}
func (c *Character) getAttackRoll() int {
	return rand.Int() % (20 + c.equipedWeapon.attack + c.Strength)
}

func (c *Character) getDefense() int {
	return c.Defense
}

func (c *Character) getStatsPage() []FormattedString {

	output := newFormattedStringCollection()
	//s1 := "LEVEL:#" + fmt.Sprintf("%2d %8s", c.level, "") + "#RACE :#" + fmt.Sprintf("%8s\n", "Human")

	output.addMessage(ct.Green, "Character Page for "+c.Name+"\n-------------------------------------------------\n")
	output.addMessage(ct.Green, "LEVEL:")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.level, ""))
	output.addMessage(ct.Green, "RACE :")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", "Human")) //TODO
	output.addMessage(ct.Green, "AGE  :")
	output.addMessage(ct.White, fmt.Sprintf("%4d %6s", 123, "")) //TODO
	output.addMessage(ct.Green, "CLASS:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", "Wizard")) //TODO
	output.addMessage(ct.Green, "STR  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Strength, ""))
	output.addMessage(ct.Green, "HitRoll:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", "66")) //TODO
	output.addMessage(ct.Green, "INT  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Inteligence, ""))
	output.addMessage(ct.Green, "DmgRoll:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", "66")) //TODO
	output.addMessage(ct.Green, "WIS  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Wisdom, ""))
	output.addMessage(ct.Green, "Alignment:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", "Paragon")) //TODO
	output.addMessage(ct.Green, "DEX  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Dexterity, ""))
	output.addMessage(ct.Green, "Armour:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", "-500")) //TODO
	output.addMessage(ct.Green, "CON  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Constitution, ""))
	output.addMessage(ct.Green, "CHA  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Charisma, ""))

	return output.fmtedStrings
}

func (c *Character) saveCharacter() {
	//TODO saveCharacter

	var ch CharacterXML
	ch.Name = c.Name
	ch.RoomIN = c.RoomIN
	ch.Defense = c.Defense
	ch.HP = c.currentHP

}

//TODO add items, stats, and any other missing fields
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
	//TODO add proper error checking, i.e. check if file exist
	xmlFile, err := os.Open("Characters/" + charName + ".xml")
	checkError(err, true)
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var charData CharacterXML
	xml.Unmarshal(XMLdata, &charData)

	char := newCharacter(charData.Name, charData.RoomIN, charData.HP, charData.Defense)

	return char
}
