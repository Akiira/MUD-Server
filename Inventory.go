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
	items         map[string]Item
	numberOfItems int
}

func newInventory() *Inventory {

	i := new(Inventory)
	i.numberOfItems = 0
	i.items = make(map[string]Item)
	return i
}

func inventoryFromXML(invXml *InventoryXML) *Inventory {
	//TODO
	ivn := new(Inventory)
	return ivn
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
	output := make([]FormattedString, len(inv.items)+1, len(inv.items)+1)

	output[0].Color = ct.White
	output[0].Value = "\nYou are carrying " + strconv.Itoa(len(inv.items)) + " unique items:"

	i := 1
	for key, _ := range inv.items {
		output[i].Color = ct.Green
		output[i].Value = "\t" + key
		i++
	}

	return output
}

func (inv *Inventory) toXML() *InventoryXML {
	return nil //TODO
}
