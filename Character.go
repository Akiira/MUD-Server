package main

import (
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math/rand"
	"net"
	"strconv"
)

type Character struct {
	Agent

	Race      string
	Class     string
	Age       int
	Alignment int

	level      int
	experience int
	gold       int

	//Equipment fields
	PersonalInvetory Inventory
	equipedWeapon    *Weapon
	equippedArmour   ArmourSet

	myClientConn *ClientConnection
}

//=================== CONSTRUCTORS =====================//

func newCharacter(name string, room int, hp int, def int) *Character {
	char := new(Character)
	char.Name = name
	char.currentHP = hp
	char.PersonalInvetory = *newInventory()
	char.equippedArmour = *newArmourSet()

	return char
}

func characterFromXML(charData *CharacterXML) *Character {
	char := new(Character)

	char.setAgentStatsFromXML(charData)

	char.gold = charData.Gold
	char.Race = charData.Race
	char.Class = charData.Class

	char.level = charData.Level
	char.experience = charData.Experience

	char.equipedWeapon = weaponFromXML(&charData.EquipedWeapon)
	char.equippedArmour = *armourSetFromXML(&charData.ArmSet)
	char.PersonalInvetory = *inventoryFromXML(&charData.PersInv)

	return char
}

//================== CLASS FUNCTIONS =============//

func (c *Character) EquipArmorByName(name string) []FormattedString {
	if item, found := c.getItemFromInv(name); found {
		if item.getItemType() == ARMOUR {
			return c.EquipArmour(item.(*Armour))
		}

		return newFormattedStringSplice("\nThat item is not armour.\n")
	}
	return newFormattedStringSplice("\nYou don't have that item. If it is on the ground try 'get'ing it first.\n")
}

func (c *Character) EquipArmour(armr *Armour) []FormattedString {
	if c.equippedArmour.isArmourEquippedAtLocation(armr.wearLocation) { // already an item present
		return newFormattedStringSplice("\nYou already have a peice of armour equiped there.\n")
	} else {
		c.equippedArmour.equipArmour(armr)
		return newFormattedStringSplice("\nYou equiped " + armr.name + "\n")
	}
}

func (c *Character) UnEquipArmourByName(name string) []FormattedString {
	armr := c.equippedArmour.takeOffArmourByName(name)
	if armr == nil {
		return newFormattedStringSplice("\nYou are not wearing a peice of armour with that name. \n")
	}
	c.addItemToInventory(armr)
	return newFormattedStringSplice("\nYou took off the " + armr.name + ".\n")
}

func (c *Character) UnWieldWeapon() []FormattedString {
	weapon := c.equipedWeapon

	if weapon == nil {
		return newFormattedStringSplice("\nYou are not wielding any weapons at this time\n")
	}
	c.addItemToInventory(weapon)
	c.equipedWeapon = nil
	return newFormattedStringSplice("\nYou put the " + weapon.name + " back in your bag.\n")
}

func (c *Character) WieldWeaponByName(name string) []FormattedString {

	if item, found := c.getItemFromInv(name); found {
		if item.getItemType() == WEAPON {
			return c.WieldWeapon(item.(*Weapon))
		}

		return newFormattedStringSplice("\nThat item is not a weapon.\n")
	}

	return newFormattedStringSplice("\nYou don't have that item. If it is on the ground try 'get'ing it first.\n")
}

func (c *Character) WieldWeapon(weapon *Weapon) []FormattedString {
	if c.equipedWeapon == nil {
		c.equipedWeapon = weapon
		return newFormattedStringSplice("\nYou equiped " + weapon.name + "\n")
	}
	return newFormattedStringSplice("\nYou already have a weapon equiped.\n")
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

func (char *Character) moveCharacter(source *Room, destination *Room) (int, []FormattedString) {

	if destination != nil {

		if destination.isLocal() {
			source.removePCFromRoom(char.Name)
			destination.addPCToRoom(char)

			return GAMEPLAY, destination.getRoomDescription()
		} else {

			source.removePCFromRoom(char.Name)
			destination.addPCToRoom(char)
			charXML := char.toXML()
			charXML.CurrentWorld = destination.WorldID
			sendCharactersXML(charXML)

			return REDIRECT, newFormattedStringSplice(servers[destination.WorldID])
		}
	} else {
		return GAMEPLAY, newFormattedStringSplice("No exit in that direction\n")
	}
}

func (char *Character) makeAttack(target Agenter) []FormattedString {

	if target == nil { // check that the target is still there, not dead from the previous round
		return newFormattedStringSplice("\nYour target does not exist any more!\n")
	}

	if !target.isDead() {

		a1 := char.getAttack()
		if a1 >= target.getDefense() {
			dmg := char.getDamage()
			target.takeDamage(dmg, 0)
			fmt.Printf("\tPlayer did %d damage.\n", char.getDamage())

			if target.isDead() {
				// TODO  reward player exp
				room := char.myClientConn.CurrentEM.worldRooms[char.RoomIN] //TODO fix this line
				room.killOffMonster(target.getName())

				return newFormattedStringSplice("You hit " + target.getName() + " for " + strconv.Itoa(dmg) + " damage and it drops over dead.\n")
			} else {
				target.addTarget(char)
				var output []FormattedString
				output = append(output, newFormattedString2(ct.Red, "\nYou were hit by "+char.Name+" for "+strconv.Itoa(dmg)+" damage!\n"))
				target.sendMessage(newServerMessageFS(output))
				return newFormattedStringSplice("\nYou hit " + target.getName() + " for " + strconv.Itoa(dmg) + " damage!\n")
			}
		} else {
			return newFormattedStringSplice("\nYour attack missed " + target.getName() + "!\n")
		}
	}

	return newFormattedStringSplice("\nthe " + target.getName() + " is already dead!\n")
}

func (c *Character) takeDamage(amount int, typeOfDamge int) {

	c.currentHP -= amount
}

func (c *Character) respawn() *FmtStrCollection {
	output := newFormattedStringCollection()
	output.addMessage(ct.Red, "\nYou died!\n")
	output.addMessage2("\nYou were respawned.\n")

	//These two lines are kinda ugly, maybe when a player dies the monster adds a respawn event to em
	// and then the even manager passes the respawn room to the characters respawn function.
	src := c.getClientConnection().CurrentEM.worldRooms[c.getClientConnection().getCharactersRoomID()]
	dest := c.getClientConnection().CurrentEM.worldRooms[worldRespawnRoomID]

	c.moveCharacter(src, dest)
	c.currentHP = c.MaxHitPoints
	return output
}
func (c *Character) isDead() bool {
	return c.currentHP <= 0 || c.myClientConn.isConnectionClosed()
}

func (c *Character) getName() string {
	return c.Name
}
func (c *Character) getRoomID() int {
	return c.RoomIN
}
func (c *Character) getAttack() int {
	return (rand.Int() % 20) + c.equipedWeapon.attack + c.Strength
}

func (c *Character) getDefense() int {
	return c.equippedArmour.getArmoursDefense()
}

func (c *Character) getClientConnection() *ClientConnection {
	return c.myClientConn
}

func (c *Character) getItemFromInv(name string) (Item_I, bool) {
	return c.PersonalInvetory.getItemByName(name)
}

func (c *Character) getAlignment() string {
	if c.Alignment > 400 {
		return "Good"
	} else if c.Alignment < -400 {
		return "Evil"
	} else {
		return "Nuetral"
	}
}

func (c *Character) getAttackRoll() string {
	min := c.equipedWeapon.attack + c.Strength
	max := 20 + c.equipedWeapon.attack + c.Strength
	return fmt.Sprintf("%d to %d", min, max)
}

func (c *Character) getDamageRoll() string {
	min := c.equipedWeapon.minDmg + c.Strength
	max := c.equipedWeapon.maxDmg + c.Strength
	return fmt.Sprintf("%d to %d", min, max)
}

func (c *Character) getDamage() int {
	return c.equipedWeapon.getDamage() + c.Strength
}

func (c *Character) getGoldAmount() int {
	return c.gold
}

func (c *Character) getStatsPage() []FormattedString {

	output := newFormattedStringCollection()

	output.addMessage(ct.Green, "Character Page for "+c.Name+"\n-------------------------------------------------\n")
	output.addMessage(ct.Green, "LEVEL:")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.level, ""))
	output.addMessage(ct.Green, "RACE :")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.Race))
	output.addMessage(ct.Green, "AGE  :")
	output.addMessage(ct.White, fmt.Sprintf("%4d %6s", c.Age, ""))
	output.addMessage(ct.Green, "CLASS:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.Class))
	output.addMessage(ct.Green, "STR  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Strength, ""))
	output.addMessage(ct.Green, "HitRoll:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.getAttackRoll()))
	output.addMessage(ct.Green, "INT  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Inteligence, ""))
	output.addMessage(ct.Green, "DmgRoll:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.getDamageRoll()))
	output.addMessage(ct.Green, "WIS  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Wisdom, ""))
	output.addMessage(ct.Green, "Alignment:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.getAlignment()))
	output.addMessage(ct.Green, "DEX  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Dexterity, ""))
	output.addMessage(ct.Green, "Armour:")
	output.addMessage(ct.White, fmt.Sprintf("%8d\n", c.getDefense()))
	output.addMessage(ct.Green, "CON  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Constitution, ""))
	output.addMessage(ct.Green, "Gold:")
	output.addMessage(ct.White, fmt.Sprintf("%8d\n", c.getGoldAmount()))
	output.addMessage(ct.Green, "CHA  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Charisma, ""))

	return output.fmtedStrings
}

func (char *Character) toXML() *CharacterXML {

	ch := new(CharacterXML)
	ch.Name = char.Name
	ch.RoomIN = char.RoomIN
	ch.HP = char.MaxHitPoints

	ch.Gold = char.gold
	ch.Class = char.Class
	ch.Race = char.Race

	ch.Strength = char.Strength
	ch.Constitution = char.Constitution
	ch.Dexterity = char.Dexterity
	ch.Wisdom = char.Wisdom
	ch.Charisma = char.Charisma
	ch.Inteligence = char.Inteligence

	ch.Level = char.level
	ch.Experience = char.experience

	ch.WeaponComment = []byte("Equipped Weapon")
	ch.EquipedWeapon = *char.equipedWeapon.toXML().(*WeaponXML)
	ch.ArmSet = *char.equippedArmour.toXML()
	ch.PersInv = *char.PersonalInvetory.toXML()

	return ch
}

func (c *Character) addTarget(target Agenter) {
	//Do nothing
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
	Race    string   `xml:"Race"`
	Class   string   `xml:"Class"`

	Strength     int `xml:"Strength"`
	Constitution int `xml:"Constitution"`
	Dexterity    int `xml:"Dexterity"`
	Wisdom       int `xml:"Wisdom"`
	Charisma     int `xml:"Charisma"`
	Inteligence  int `xml:"Inteligence"`

	Level      int `xml:"Level"`
	Experience int `xml:"Experience"`
	Gold       int `xml:"Gold"`

	CurrentWorld string `xml:"CurrentWorld"`

	WeaponComment xml.Comment  `xml:",comment"`
	EquipedWeapon WeaponXML    `xml:"Weapon"`
	ArmSet        ArmourSetXML `xml:"ArmourSet"`
	PersInv       InventoryXML `xml:"Inventory"`
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

	char := characterFromXML(&queriedChar)

	return char
}

//func saveCharacterToFile(char *Character) {
//	//TODO saveCharacter

//	var ch CharacterXML
//	ch.Name = char.Name
//	ch.RoomIN = char.RoomIN
//	ch.HP = char.currentHP

//}
