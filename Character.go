package main

import (
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
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
	char.equippedArmour = *NewArmourSet()

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
	char.equippedArmour = *NewArmourSetFromXML(&charData.ArmSet)
	char.PersonalInvetory = *inventoryFromXML(&charData.PersInv)

	return char
}

//============================== CLASS FUNCTIONS =============================//

// ===== Armour Functions

func (c *Character) EquipArmorByName(name string) []FormattedString {
	if item, found := c.GetItem(name); found {
		if item.getItemType() == ARMOUR {
			return c.EquipArmour(item.(*Armour))
		}

		return newFormattedStringSplice("\nThat item is not armour.\n")
	}
	return newFormattedStringSplice("\nYou don't have that item. If it is on the ground try 'get'ing it first.\n")
}

func (c *Character) EquipArmour(armr *Armour) []FormattedString {
	if c.equippedArmour.IsArmourAt(armr.wearLocation) { // already an item present
		return newFormattedStringSplice("\nYou already have a peice of armour equiped there.\n")
	} else {
		c.equippedArmour.EquipArmour(armr)
		return newFormattedStringSplice("\nYou equiped " + armr.name + "\n")
	}
}

func (c *Character) UnEquipArmourByName(name string) []FormattedString {
	armr := c.equippedArmour.GetAndRemoveArmour(name)
	if armr == nil {
		return newFormattedStringSplice("\nYou are not wearing a peice of armour with that name. \n")
	}
	c.AddItemToInventory(armr)
	return newFormattedStringSplice("\nYou took off the " + armr.name + ".\n")
}

func (c *Character) UnEquipArmourAt(location string) []FormattedString {
	if c.equippedArmour.IsArmourAt(location) {
		armr := c.equippedArmour.GetAndRemoveArmour(location)
		c.AddItemToInventory(armr)
		return newFormattedStringSplice(fmt.Sprintf("You succesfully removed the %s and stored it in your inventory.\n", armr.getName()))
	} else {
		return newFormattedStringSplice("You are not wearing any armour there.\n")
	}
}

// ===== Weapon Functions

func (c *Character) WieldWeapon(weapon interface{}) []FormattedString {
	var weaponToEquip Item_I

	switch weapon := weapon.(type) {
	default:
		fmt.Printf("Unexpected type %T in Character.WieldWeapon", weapon)
	case string:
		if item, found := c.GetItem(weapon); found {
			weaponToEquip = item
		}
	case *Weapon:
		weaponToEquip = weapon
	}

	if weaponToEquip == nil {
		return newFormattedStringSplice("\nYou don't have that item. If it is on the ground try 'get'ing it first.\n")
	} else if weaponToEquip.getItemType() != WEAPON {
		return newFormattedStringSplice("\nThat item is not a weapon.\n")
	} else if c.equipedWeapon == nil {
		c.equipedWeapon = weaponToEquip.(*Weapon)
		return newFormattedStringSplice("\nYou equiped " + c.equipedWeapon.getName() + ".\n")
	} else {
		return newFormattedStringSplice("\nYou already have a weapon equiped.\n")
	}
}

func (c *Character) UnWieldWeapon() []FormattedString {
	if weapon := c.equipedWeapon; weapon != nil {
		c.AddItemToInventory(weapon)
		c.equipedWeapon = nil
		return newFormattedStringSplice("\nYou put the " + weapon.name + " back in your bag.\n")
	} else {
		return newFormattedStringSplice("\nYou are not wielding any weapons at this time\n")
	}
}

// ===== Inventory Functions

func (c *Character) AddInventoryToInventory(otherInv *Inventory) {
	c.PersonalInvetory.AddInventory(otherInv)
}

func (c *Character) AddItemsToInventory(items []Item_I) {
	c.PersonalInvetory.AddItems(items)
}

func (c *Character) AddItemToInventory(item Item_I) {
	c.PersonalInvetory.AddItem(item)
}

// ===== General Functions

func (char *Character) moveCharacter(source *Room, destination *Room) (int, []FormattedString) {

	if destination != nil {

		if destination.IsLocal() {
			source.RemovePlayer(char.Name)
			destination.AddPlayer(char)

			return GAMEPLAY, destination.GetDescription()
		} else {
			source.RemovePlayer(char.Name)
			destination.AddPlayer(char)

			charXML := char.toXML()
			charXML.CurrentWorld = destination.WorldID

			SendCharactersXML(charXML)

			return REDIRECT, newFormattedStringSplice(servers[destination.WorldID])
		}
	} else {
		return GAMEPLAY, newFormattedStringSplice("No exit in that direction\n")
	}
}

func (char *Character) makeAttack(target Agenter) []FormattedString {

	if target == nil { // check that the target is still there, not dead from the previous round
		return newFormattedStringSplice("\nYour target does not exist, perhaps you typed the name wrong or another player killed it!\n")
	}

	if !target.IsDead() {

		a1 := char.GetAttack()
		if a1 >= target.GetDefense() {
			dmg := char.getDamage()
			target.takeDamage(dmg, 0)
			target.takeDamage(char.getDamage(), 0)
			//fmt.Printf("\tPlayer did %d damage.\n", char.getDamage())

			if target.IsDead() {
				// TODO  reward player exp
				room := char.myClientConn.EventManager.worldRooms[char.RoomIN] //TODO fix this line
				room.killOffMonster(target.GetName())

				return newFormattedStringSplice("You hit " + target.GetName() + " for " + strconv.Itoa(dmg) + " damage and it drops over dead.\n")
			} else {
				target.addTarget(char)
				var output []FormattedString
				output = append(output, newFormattedString2(ct.Red, "\nYou were hit by "+char.Name+" for "+strconv.Itoa(dmg)+" damage!\n"))
				target.SendMessage(newServerMessageFS(output))
				return newFormattedStringSplice("\nYou hit " + target.GetName() + " for " + strconv.Itoa(dmg) + " damage!\n")
			}
		} else {

			return newFormattedStringSplice("\nYour attack missed " + target.GetName() + "!\n")
		}
	}

	return newFormattedStringSplice("\nthe " + target.GetName() + " is already dead!\n")
}

func (c *Character) takeDamage(amount int, typeOfDamge int) {
	c.currentHP -= amount
}

func (c *Character) addTarget(target Agenter) {
	//Do nothing, required for Agenter interface
}

func (c *Character) ApplyFleePenalty() []FormattedString {
	c.experience -= 100 * c.level
	return newFormattedStringSplice2(ct.Red, fmt.Sprintf("\nYou lost %d experience for fleeing.\n", 100*c.level))
}

func (c *Character) respawn() *FmtStrCollection {
	output := newFormattedStringCollection()
	output.addMessage(ct.Red, "\nYou died!\n")
	output.addMessage2("\nYou were respawned.\n")

	//These two lines are kinda ugly, maybe when a player dies the monster adds a respawn event to em
	// and then the even manager passes the respawn room to the characters respawn function.
	src := c.GetClientConnection().EventManager.worldRooms[c.GetClientConnection().getCharactersRoomID()]
	dest := c.GetClientConnection().EventManager.worldRooms[worldRespawnRoomID]

	c.moveCharacter(src, dest)
	c.currentHP = c.MaxHitPoints
	return output
}

// ===== Getter Functions

func (c *Character) GetName() string {
	return c.Name
}

func (c *Character) GetRoomID() int {
	return c.RoomIN
}

func (c *Character) GetAttack() int {
	return (rand.Int() % 20) + c.equipedWeapon.attack + c.Strength
}

func (c *Character) GetDefense() int {
	return c.equippedArmour.GetDefense()
}

//TODO refactor so we can remove this. The character class should not provide
//open access to its client connection.
func (c *Character) GetClientConnection() *ClientConnection {
	return c.myClientConn
}

func (c *Character) GetAndRemoveItems(names []string) (items []Item_I) {
	for _, name := range names {
		if item, found := c.GetAndRemoveItem(name); found {
			items = append(items, item)
		}
	}

	return items
}

func (c *Character) GetAndRemoveItem(name string) (Item_I, bool) {
	return c.PersonalInvetory.GetAndRemoveItem(name)
}

func (c *Character) GetItem(name string) (Item_I, bool) {
	return c.PersonalInvetory.GetItem(name)
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
	return c.equipedWeapon.GetDamage() + c.Strength
}

func (c *Character) GetGoldAmount() int {
	return c.gold
}

func (c *Character) GetItemsToTrade(inv *Inventory, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		response := c.myClientConn.GetItemsToTrade()

		if response == "timeout" || response == "done" {
			break
		} else {
			if item, found := c.GetAndRemoveItem(response); found {
				inv.AddItem(item)
				c.SendMessage("One " + response + " was added to the trade pool.\n")
			} else {
				c.SendMessage("You do not have any more of the item: " + response + ".\n")
			}
		}
	}
}

func (c *Character) GetTradeResponse() string {
	return c.myClientConn.GetResponseToTrade()
}

func (c *Character) GetEquipment() []FormattedString {
	output := newFormattedStringCollection()

	output.addMessage(ct.Green, "\t\t\tEquipment Page for "+c.Name+"\n--------------------------------------------------------------------\n")
	if c.equipedWeapon != nil {
		output.addMessages2(c.equipedWeapon.GetWeaponPage())
	} else {
		output.addMessage(ct.Green, "\t\t\tEquipped Weapon\n")
		output.addMessage2(fmt.Sprintf("\t%-15s   %-15s %-15s %-15s\n", "Name", "Attack", "MinDmg", "MaxDmg"))
		output.addMessage(ct.Green, "--------------------------------------------------------------------\n\n")
	}
	output.addMessages2(c.equippedArmour.GetArmourWornPage())

	return output.fmtedStrings
}

func (c *Character) GetStats() []FormattedString {

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
	output.addMessage(ct.White, fmt.Sprintf("%8d\n", c.GetDefense()))
	output.addMessage(ct.Green, "CON  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Constitution, ""))
	output.addMessage(ct.Green, "Gold:")
	output.addMessage(ct.White, fmt.Sprintf("%8d\n", c.GetGoldAmount()))
	output.addMessage(ct.Green, "CHA  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Charisma, ""))

	return output.fmtedStrings
}

// ===== Predicate Functions

func (c *Character) HasItems(itemNames []string) bool {
	for _, name := range itemNames {
		if !c.HasItem(name) {
			return false
		}
	}

	return true
}

func (c *Character) HasItem(name string) bool {
	return c.PersonalInvetory.PossesItem(name)
}

func (c *Character) IsDead() bool {
	return c.currentHP <= 0 || c.myClientConn.isConnectionClosed()
}

// ===== Misc Functions

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

	ch.CurrentWorld = os.Args[1]

	return ch
}

//SendMessage is used to send a communication to a character. This communication
//is then forwarded to the client for this character. The only valid message types
//are: string, ServerMessage, and []FormattedString.
func (char *Character) SendMessage(msg interface{}) {
	switch msg := msg.(type) {
	default:
		fmt.Printf("Unexpected type %T in Character.SendMessage\n", msg)
	case string:
		char.myClientConn.Write(newServerMessageS(msg))
	case ServerMessage:
		char.myClientConn.Write(msg)
	case []FormattedString:
		char.myClientConn.Write(newServerMessageFS(msg))
	}
}

//==============="STATIC" FUNCTIONS===================//

//TODO add any missing fields
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

//GetCharacterFromStorage querys the characterStorage server for the characters
//data and stores it into a new character object.
func GetCharacterFromStorage(charName string) (char *Character, err error) {
	conn, err := net.Dial("tcp", servers["characterStorage"])
	if err != nil {
		return char, err
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	var queriedChar CharacterXML

	//Request file from server
	if err = enc.Encode(newServerMessageTypeS(GETFILE, charName)); err != nil {
		return char, err
	}

	//Decone the character from message
	if err = dec.Decode(&queriedChar); err != nil {
		return char, err
	}

	//Decode the characters xml into a character object
	char = characterFromXML(&queriedChar)
	return char, err
}
