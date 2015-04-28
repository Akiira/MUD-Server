package main

import (
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math/rand"
	"net"
	"strconv"
	"sync"
)

type Character struct {
	Agent

	Race      string
	Class     string
	Age       int
	Alignment int

	experience int
	gold       int

	//Equipment fields
	PersonalInvetory Inventory
	equipedWeapon    *Weapon
	equippedArmour   ArmourSet

	myClientConn *ClientConnection
	inv_mutex    sync.Mutex
}

//=================== CONSTRUCTOR ======================//

func NewCharacter(charData *CharacterXML) *Character {
	char := new(Character)

	char.SetAgentStats(charData)

	char.gold = charData.Gold
	char.Race = charData.Race
	char.Class = charData.Class

	char.experience = charData.Experience

	char.equipedWeapon = NewWeaponFromXML(&charData.EquipedWeapon)
	char.equippedArmour = *NewArmourSetFromXML(&charData.ArmSet)
	char.PersonalInvetory = *NewInventoryFromXML(&charData.PersInv)
	char.inv_mutex = sync.Mutex{}

	return char
}

//============================== CLASS FUNCTIONS =============================//

// ===== Armour Functions

func (c *Character) EquipArmor(name string) []FormattedString {
	if item, found := c.GetItem(name); found {
		if item.GetType() == ARMOUR {
			return c.equipArmour(item.(*Armour))
		}

		return newFormattedStringSplice("\nThat item is not armour.\n")
	}
	return newFormattedStringSplice("\nYou don't have that item. If it is on the ground try 'get'ing it first.\n")
}

func (c *Character) equipArmour(armr *Armour) []FormattedString {
	if c.equippedArmour.IsArmourAt(armr.wearLocation) { // already an item present
		return newFormattedStringSplice("\nYou already have a peice of armour equiped there.\n")
	} else {
		c.equippedArmour.EquipArmour(armr)
		return newFormattedStringSplice("\nYou equiped " + armr.name + "\n")
	}
}

func (c *Character) UnEquipArmour(name string) []FormattedString {
	armr := c.equippedArmour.GetAndRemoveArmour(name)
	if armr == nil {
		return newFormattedStringSplice("\nYou are not wearing a peice of armour with that name. \n")
	}
	c.AddItem(armr)
	return newFormattedStringSplice("\nYou took off the " + armr.name + ".\n")
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
	} else if weaponToEquip.GetType() != WEAPON {
		return newFormattedStringSplice("\nThat item is not a weapon.\n")
	} else if c.equipedWeapon == nil {
		c.equipedWeapon = weaponToEquip.(*Weapon)
		return newFormattedStringSplice("\nYou equiped " + c.equipedWeapon.GetName() + ".\n")
	} else {
		return newFormattedStringSplice("\nYou already have a weapon equiped.\n")
	}
}

func (c *Character) UnWieldWeapon() []FormattedString {
	if weapon := c.equipedWeapon; weapon != nil {
		c.AddItem(weapon)
		c.equipedWeapon = nil
		return newFormattedStringSplice("\nYou put the " + weapon.name + " back in your bag.\n")
	} else {
		return newFormattedStringSplice("\nYou are not wielding any weapons at this time\n")
	}
}

// ===== Inventory Functions

func (c *Character) AddInventory(otherInv *Inventory) {
	c.PersonalInvetory.AddInventory(otherInv)
}

func (c *Character) AddItems(items []Item_I) {
	c.PersonalInvetory.AddItems(items)
}

func (c *Character) AddItem(item Item_I) {
	if item.GetName() == "Gold" {
		c.gold += item.GetWorth()
	} else {
		c.PersonalInvetory.AddItem(item)
	}
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
	c.inv_mutex.Lock()
	defer c.inv_mutex.Unlock()

	return c.PersonalInvetory.GetAndRemoveItem(name)
}

func (c *Character) GetItem(name string) (Item_I, bool) {
	return c.PersonalInvetory.GetItem(name)
}

// ===== General Functions

func (char *Character) Move(source *Room, destination *Room) (int, []FormattedString) {

	if source != nil && destination != nil {
		source.RemovePlayer(char.Name)

		if destination.IsLocal() {
			destination.AddPlayer(char)

			return GAMEPLAY, destination.GetDescription()
		} else {
			destination.AddPlayer(char)
			SendCharactersXML(char.ToXML())

			return REDIRECT, newFormattedStringSplice(servers[destination.WorldID])
		}
	} else {
		return GAMEPLAY, newFormattedStringSplice("No exit in that direction\n")
	}
}

//Attack executes an attack, by the character, against the supplied target. It is
//required for the Agenter interface. Player experience is automatically handled
//here. If the target dies the correct clean up functions are called.
func (char *Character) Attack(target Agenter) []FormattedString {

	if target == nil {
		return newFormattedStringSplice("\nYour target does not exist, perhaps you typed the name wrong or another player killed it!\n")
	}

	if !target.IsDead() {

		a1 := char.GetAttack()
		if a1 >= target.GetDefense() {
			dmg := char.GetDamage()
			target.TakeDamage(dmg, 0)
			target.TakeDamage(char.GetDamage(), 0)

			if target.IsDead() {
				char.experience += 10 * target.GetLevel()
				room := eventManager.GetRoom(char)
				room.KillOffMonster(target.GetName())

				return newFormattedStringSplice("You hit " + target.GetName() + " for " + strconv.Itoa(dmg) + " damage and it drops over dead.\n")
			} else {
				target.AddTarget(char)
				char.experience += dmg / 2

				var output []FormattedString
				output = append(output, newFormattedString2(ct.Red, "\nYou were hit by "+char.Name+" for "+strconv.Itoa(dmg)+" damage!\n"))
				target.SendMessage(newServerMessageFS(output))
				return newFormattedStringSplice("\nYou hit " + target.GetName() + " for " + strconv.Itoa(dmg) + " damage!\n")
			}
		} else {
			char.experience += 2
			return newFormattedStringSplice("\nYour attack missed " + target.GetName() + "!\n")
		}
	}

	return newFormattedStringSplice("\nthe " + target.GetName() + " is already dead!\n")
}

func (c *Character) TakeDamage(amount int, typeOfDamge int) {
	c.currentHP -= amount
}

func (c *Character) AddTarget(target Agenter) {
	//Do nothing, required for Agenter interface
}

func (c *Character) ApplyFleePenalty() []FormattedString {
	c.experience -= 100 * c.Level
	return newFormattedStringSplice2(ct.Red, fmt.Sprintf("\nYou lost %d experience for fleeing.\n", 100*c.Level))
}

func (c *Character) Respawn() []FormattedString {
	output := newFormattedStringCollection()
	output.addMessage(ct.Red, "\nYou died!\n")
	output.addMessage2("\nYou repspawned.\n")

	src := eventManager.GetRoom(c)
	dest := eventManager.GetRespawnRoom()

	c.Move(src, dest)
	c.currentHP = c.MaxHitPoints
	return output.fmtedStrings
}

func (c *Character) LevelUp() []FormattedString {
	if c.experience-(1000*c.GetLevel()) > 0 {
		c.incrementAttack()
		c.incrementDefense()
		c.incrementLevel()

		return newFormattedStringSplice2(ct.Green, "\nYou leveled up!\n")
	} else {
		return newFormattedStringSplice2(ct.Green, "\nYou don't have enough exp to level up yet.\n")
	}
}

// ===== Getter Functions

func (c *Character) GetName() string {
	return c.Name
}

func (c *Character) GetDescription() string {
	return c.Name //TODO GetDescription
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

func (c *Character) GetAlignment() string {
	if c.Alignment > 400 {
		return "Good"
	} else if c.Alignment < -400 {
		return "Evil"
	} else {
		return "Nuetral"
	}
}

func (c *Character) GetAttackRoll() string {
	min := c.equipedWeapon.attack + c.Strength
	max := 20 + c.equipedWeapon.attack + c.Strength
	return fmt.Sprintf("%d to %d", min, max)
}

func (c *Character) GetDamageRoll() string {
	min := c.equipedWeapon.minDmg + c.Strength
	max := c.equipedWeapon.maxDmg + c.Strength
	return fmt.Sprintf("%d to %d", min, max)
}

func (c *Character) GetDamage() int {
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
				c.SendMessage("\nOne " + item.GetName() + " was added to the trade pool.\n")
			} else {
				c.SendMessage("\nYou do not have any more of the item: " + response + ".\n")
			}
		}
	}
}

func (c *Character) GetTradeResponse() string {
	return c.myClientConn.GetResponseToTrade()
}

func (c *Character) GetEquipmentPage() []FormattedString {
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

func (c *Character) GetStatsPage() []FormattedString {

	output := newFormattedStringCollection()

	output.addMessage(ct.Green, "Character Page for "+c.Name+"\n-------------------------------------------------\n")
	output.addMessage(ct.Green, "LEVEL:")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Level, ""))
	output.addMessage(ct.Green, "RACE :")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.Race))
	output.addMessage(ct.Green, "AGE  :")
	output.addMessage(ct.White, fmt.Sprintf("%4d %6s", c.Age, ""))
	output.addMessage(ct.Green, "CLASS:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.Class))
	output.addMessage(ct.Green, "STR  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Strength, ""))
	output.addMessage(ct.Green, "HitRoll:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.GetAttackRoll()))
	output.addMessage(ct.Green, "INT  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Inteligence, ""))
	output.addMessage(ct.Green, "DmgRoll:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.GetDamageRoll()))
	output.addMessage(ct.Green, "WIS  :")
	output.addMessage(ct.White, fmt.Sprintf("%2d %8s", c.Wisdom, ""))
	output.addMessage(ct.Green, "Alignment:")
	output.addMessage(ct.White, fmt.Sprintf("%8s\n", c.GetAlignment()))
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
	return c.currentHP <= 0 || c.myClientConn.IsConnectionClosed()
}

// ===== Misc Functions

func (char *Character) ToXML() *CharacterXML {

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

	ch.Level = char.Level
	ch.Experience = char.experience

	ch.WeaponComment = []byte("Equipped Weapon")
	ch.EquipedWeapon = *char.equipedWeapon.ToXML().(*WeaponXML)
	ch.ArmSet = *char.equippedArmour.ToXML()
	ch.PersInv = *char.PersonalInvetory.toXML()

	ch.CurrentWorld = eventManager.GetPlayersWorld(char)

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

	//Decode the character from message
	if err = dec.Decode(&queriedChar); err != nil {
		return char, err
	}

	//Decode the characters xml into a character object
	char = NewCharacter(&queriedChar)
	return char, err
}
