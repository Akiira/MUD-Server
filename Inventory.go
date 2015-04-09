package main

import (
	"encoding/xml"
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
	items   map[string]Item
	weapons map[string]Weapon
	armours map[string]Armour
}

func newInvXML() *InventoryXML {
	invXML := new(InventoryXML)
	invXML.Items = make([]ItemXML)
	invXML.Weapons = make([]WeaponXML)
	invXML.Armours = make([]ArmourXML)

	return invXML
}

func newInventory() *Inventory {

	i := new(Inventory)
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
	}
}
func (inv *Inventory) addWeaponToInventory(addWeapon *Weapon) {
	if val, ok := inv.weapons[addWeapon.Item.name]; ok { // the item is already there
		val.quantity++
		inv.weapons[addWeapon.Item.name] = val
	} else {
		inv.weapons[addWeapon.Item.name] = *addWeapon
	}
}
func (inv *Inventory) addArmourToInventory(addArmour *Armour) {
	if val, ok := inv.armours[addArmour.Item.name]; ok { // the item is already there
		val.quantity++
		inv.armours[addArmour.Item.name] = val
	} else {
		inv.armours[addArmour.Item.name] = *addArmour
	}
}

func inventoryFromXML(invXml *InventoryXML) *Inventory {
	inv := newInventory()

	//loop through items
	for _, itm := range invXml.Items {
		inv.addItemToInventory(itemFromXML(&itm))
	}

	//loop through weapons
	for _, itm := range invXml.Weapons {
		inv.addWeaponToInventory(weaponFromXML(&itm))
	}

	//loop through armours
	for _, itm := range invXml.Armours {
		inv.addArmourToInventory(armourFromXML(&itm))
	}

	return inv
}

func (inv *Inventory) getItemByName(name string) *Item {
	itm, _ := inv.items[name]
	return &itm
}

func (inv *Inventory) getInventoryDescription() []FormattedString {
	output := make([]FormattedString, len(inv.items)+len(inv.weapons)+len(inv.armours)+1, len(inv.items)+len(inv.weapons)+len(inv.armours)+1)

	output[0].Color = ct.White
	output[0].Value = "\nYou are carrying " + strconv.Itoa(len(inv.items)+len(inv.weapons)+len(inv.armours)) + " unique items:\n"

	i := 1
	for key, item := range inv.items {
		output[i].Color = ct.Green
		output[i].Value = "\t" + strconv.Itoa(item.quantity) + "\t" + key + "\n"
		i++
	}

	for key, weapon := range inv.weapons {
		output[i].Color = ct.Green
		output[i].Value = "\t" + strconv.Itoa(weapon.Item.quantity) + "\t" + key + "\n"
		i++
	}

	for key, armour := range inv.armours {
		output[i].Color = ct.Green
		output[i].Value = "\t" + strconv.Itoa(armour.Item.quantity) + "\t" + key + "\n"
		i++
	}

	return output
}

func (inv *Inventory) toXML() *InventoryXML {

	for _, item := range inv.items {
		invXML.Items = append(invXML.Items, *item.toXML())
	}

	for _, weapon := range inv.weapons {
		invXML.Weapons = append(invXML.Weapons, *weapon.toXML())
	}

	for _, armour := range inv.armours {
		invXML.Armours = append(invXML.Armours, *armour.toXML())
	}

	return invXML
}
