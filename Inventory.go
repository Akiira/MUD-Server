package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"strconv"
)

type InventoryXML struct {
	XMLName xml.Name    `xml:"Inventory"`
	Items   []ItemXML   `xml:"Item"`
	Weapons []WeaponXML `xml:"Weapon"`
	Armours []ArmourXML `xml:"Armour"`
}

type Inventory struct {
	items           map[string]Item
	weapons         map[string]Weapon
	armours         map[string]Armour
	numberOfItems   int
	numberOfWeapons int
	numberOfArmours int
}

func newInventory() *Inventory {

	i := new(Inventory)
	i.numberOfItems = 0
	i.numberOfWeapons = 0
	i.numberOfArmours = 0
	i.items = make(map[string]Item)
	i.weapons = make(map[string]Weapon)
	i.armours = make(map[string]Armour)

	return i
}

func (inv *Inventory) addItemToInventory(addItem *Item) {
	if val, ok := inv.items[addItem.name]; ok { // the item is already there
		val.quantity++
		inv.items[addItem.name] = val
	} else {
		inv.items[addItem.name] = *addItem
		inv.numberOfItems++
	}
}
func (inv *Inventory) addWeaponToInventory(addWeapon *Weapon) {
	if val, ok := inv.weapons[addWeapon.Item.name]; ok { // the item is already there
		val.quantity++
		inv.weapons[addWeapon.Item.name] = val
	} else {
		inv.weapons[addWeapon.Item.name] = *addWeapon
		inv.numberOfWeapons++
	}
}
func (inv *Inventory) addArmourToInventory(addArmour *Armour) {
	if val, ok := inv.armours[addArmour.Item.name]; ok { // the item is already there
		val.quantity++
		inv.armours[addArmour.Item.name] = val
	} else {
		inv.armours[addArmour.Item.name] = *addArmour
		inv.numberOfArmours++
	}
}

func inventoryFromXML(invXml *InventoryXML) *Inventory {
	//TODO
	inv := newInventory()
	fmt.Println("this is invXML")
	fmt.Println(invXml)

	//loop through items
	for i := 0; i < len(invXml.Items); i++ {

		if val, ok := inv.items[invXml.Items[i].Name]; ok { // the item is already there
			val.quantity++
			inv.items[invXml.Items[i].Name] = val
		} else {
			inv.items[invXml.Items[i].Name] = *itemFromXML(&(invXml.Items[i]))
			inv.numberOfItems++
		}

	}

	//loop through weapons
	for i := 0; i < len(invXml.Weapons); i++ {

		if val, ok := inv.weapons[invXml.Weapons[i].ItemInfo.Name]; ok { // the item is already there
			val.quantity++
			inv.weapons[invXml.Weapons[i].ItemInfo.Name] = val
		} else {
			inv.weapons[invXml.Weapons[i].ItemInfo.Name] = *weaponFromXML(&(invXml.Weapons[i]))
			inv.numberOfWeapons++
		}

	}

	//loop through armours
	for i := 0; i < len(invXml.Armours); i++ {

		if val, ok := inv.armours[invXml.Armours[i].ItemInfo.Name]; ok { // the item is already there
			val.quantity++
			inv.armours[invXml.Armours[i].ItemInfo.Name] = val
		} else {
			inv.armours[invXml.Armours[i].ItemInfo.Name] = *armourFromXML(&(invXml.Armours[i]))
			inv.numberOfArmours++
		}

	}

	fmt.Println(inv)

	return inv
}

//func (inv *Inventory) getItemByName(name string) Item {
//	var i int
//	for i = 0; i < inv.numberOfItems; i++ {
//		if inv.items[i].name == name {
//			return inv.items[i]
//		}
//	}
//	var null Item
//	return null
//}

func (inv *Inventory) getInventoryDescription() []FormattedString {
	output := make([]FormattedString, len(inv.items)+len(inv.weapons)+len(inv.armours)+1, len(inv.items)+len(inv.weapons)+len(inv.armours)+1)

	output[0].Color = ct.White
	output[0].Value = "\nYou are carrying " + strconv.Itoa(len(inv.items)+len(inv.weapons)+len(inv.armours)) + " unique items:"

	i := 1
	for key, item := range inv.items {
		output[i].Color = ct.Green
		output[i].Value = "\t" + strconv.Itoa(item.quantity) + "\t" + key
		i++
	}

	for key, weapon := range inv.weapons {
		output[i].Color = ct.Green
		output[i].Value = "\t" + strconv.Itoa(weapon.Item.quantity) + "\t" + key
		i++
	}

	for key, armour := range inv.armours {
		output[i].Color = ct.Green
		output[i].Value = "\t" + strconv.Itoa(armour.Item.quantity) + "\t" + key
		i++
	}

	return output
}

func (inv *Inventory) toXML() *InventoryXML {
	invXML := new(InventoryXML)

	invXML.Items = make([]ItemXML, len(inv.items), len(inv.items))

	i := 0
	for _, item := range inv.items {
		invXML.Items[i] = *(item.toXML())
		i++
	}

	invXML.Weapons = make([]WeaponXML, len(inv.weapons), len(inv.weapons))

	i = 0
	for _, weapon := range inv.weapons {
		invXML.Weapons[i] = *(weapon.toXML())
		i++
	}

	invXML.Armours = make([]ArmourXML, len(inv.armours), len(inv.armours))

	i = 0
	for _, armour := range inv.armours {
		invXML.Armours[i] = *(armour.toXML())
		i++
	}

	return invXML //TODO
}
