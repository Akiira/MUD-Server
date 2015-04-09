package main

import (
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math/rand"
	"net"
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

//=================== CONSTRUCTORS =====================//

func newCharacter(name string, room int, hp int, def int) *Character {
	char := new(Character)
	char.Name = name
	char.currentHP = hp
	char.Defense = def
	char.PersonalInvetory = *newInventory()
	char.equippedArmour = *newArmourSet()

	return char
}

func characterFromXML(charData *CharacterXML) *Character {
	char := new(Character)
	char.setAgentStatsFromXML(charData)
	char.equipedWeapon = *weaponFromXML(&charData.EquipedWeapon)
	char.equippedArmour = *armourSetFromXML(charData.ArmSet)
	char.level = charData.Level
	char.experience = charData.experience

	char.PersonalInvetory = *inventoryFromXML(&charData.PersInv)

	//TODO still not finished

	return char
}

//================== CLASS FUNCTIONS =============//

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

func (c *Character) addItemToInventory(item Item_I) {
	c.PersonalInvetory.addItemToInventory(item)
}

//TODO
//func (char *Character) moveCharacter(direction string, source *Room, destination *Room) []FormattedString

func (char *Character) moveCharacter(direction string) (int, []FormattedString) {

	dirAsInt := convertDirectionToInt(direction)
	room := char.myClientConn.CurrentEM.worldRooms[char.RoomIN] //TODO this is just a temporary fix

	if room.isValidDirection(dirAsInt) {
		newRoom := room.getConnectedRoom(dirAsInt)

		if newRoom.isLocal() {
			room.removePCFromRoom(char.Name)
			newRoom.addPCToRoom(char)

			return GAMEPLAY, room.ExitLinksToRooms[dirAsInt].getRoomDescription()
		} else {

			room.removePCFromRoom(char.Name)
			newRoom.addPCToRoom(char)

			sendCharactersXML(char.toXML())

			return REDIRECT, newFormattedStringSplice(servers[newRoom.WorldID])
		}
	} else {
		return GAMEPLAY, newFormattedStringSplice("No exit in that direction\n")
	}
}

func (char *Character) makeAttack(target Agenter) []FormattedString {

	output := make([]FormattedString, 2, 2)

	if target != nil { // check that the target is still there, not dead from the previous round
		return newFormattedStringSplice("\nYour target is not exist any more!\n")
	}

	if !target.isDead() {

		a1 := char.getAttackRoll()
		if a1 >= target.getDefense() {
			target.takeDamage(char.getDamage(), 0)
			output[0].Value = "\nYou hit the " + target.getName() + "!\n"
		} else {
			output[0].Value = "\nYou missed the " + target.getName() + "!\n"
		}

		if target.isDead() {
			// TODO  reward player exp
			output[1].Value = "The " + target.getName() + " drops over dead.\n"
			room := char.myClientConn.CurrentEM.worldRooms[char.RoomIN] //TODO fix this line
			room.killOffMonster(target.getName())
		}

	} else {
		output[0].Value = "\nthe " + target.getName() + " is already dead!\n"
	}

	return output

}

func (c *Character) takeDamage(amount int, typeOfDamge int) []FormattedString {
	//	if amount-c.equippedArmour.getArmoursDefense() > 0 {
	//		c.currentHP -= (amount - c.equippedArmour.getArmoursDefense())
	//	}
	c.currentHP -= amount
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
	return (rand.Int() % 20) + c.equipedWeapon.attack + c.Strength
}

func (c *Character) getDefense() int {
	return c.Defense
}

func (c *Character) getClientConnection() *ClientConnection {
	return c.myClientConn
}

func (c *Character) getItemFromInv(name string) (Item_I, bool) {
	return c.PersonalInvetory.getItemByName(name)
}

func (c *Character) getDamage() int {
	return c.equipedWeapon.getDamage() + c.Strength
}

func (c *Character) getStatsPage() []FormattedString {

	output := newFormattedStringCollection()

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

func (char *Character) toXML() *CharacterXML {

	ch := new(CharacterXML)
	ch.Name = char.Name
	ch.RoomIN = char.RoomIN
	ch.Defense = char.Defense
	ch.HP = char.currentHP

	ch.Strength = char.Strength
	ch.Constitution = char.Constitution
	ch.Dexterity = char.Dexterity
	ch.Wisdom = char.Wisdom
	ch.Charisma = char.Charisma
	ch.Inteligence = char.Inteligence

	ch.Level = char.level
	ch.experience = char.experience

	ch.EquipedWeapon = *char.equipedWeapon.toXML()
	ch.ArmSet = *char.equippedArmour.toXML()
	ch.PersInv = *char.PersonalInvetory.toXML()

	ch.Wpn = char.equipedWeapon.toXML()

	return ch
}

func (char *Character) sendMessage(msg ServerMessage) {
	char.myClientConn.sendMsgToClient(msg)
}

//==============="STATIC" FUNCTIONS===================//

//TODO add items, stats, and any other missing fields
type CharacterXML struct {
	XMLName xml.Name `xml:"Character"`
	Name    string   `xml:"Name"`
	RoomIN  int      `xml:"RoomIN"`
	HP      int      `xml:"HitPoints"`
	Defense int      `xml:"Defense"`

	Strength     int `xml:"Strength"`
	Constitution int `xml:"Constitution"`
	Dexterity    int `xml:"Dexterity"`
	Wisdom       int `xml:"Wisdom"`
	Charisma     int `xml:"Charisma"`
	Inteligence  int `xml:"Inteligence"`

	Level      int `xml:"Level"`
	experience int `xml:"experience"`

	CurrentWorld string `xml:"CurrentWorld"`

	EquipedWeapon WeaponXML    `xml:"Weapon"`
	ArmSet        ArmourSetXML `xml:"ArmourSet"`
	PersInv       InventoryXML `xml:"Inventory"`
	Wpn           ItemXML_I    //Testing
}

func getCharacterFromCentral(charName string) *Character {

	address := servers["characterStorage"]

	conn, err := net.Dial("tcp", address)
	checkError(err, true)
	defer conn.Close()
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	serverMsg := newServerMessageTypeS(GETFILE, charName)

	var queriedChar CharacterXML

	err = enc.Encode(serverMsg)
	checkError(err, true)
	err = dec.Decode(&queriedChar)
	checkError(err, true)

	fmt.Println("this is received char")
	fmt.Println(queriedChar)

	char := characterFromXML(&queriedChar)

	return char
}

func saveCharacterToFile(char *Character) {
	//TODO saveCharacter

	var ch CharacterXML
	ch.Name = char.Name
	ch.RoomIN = char.RoomIN
	ch.Defense = char.Defense
	ch.HP = char.currentHP

}
