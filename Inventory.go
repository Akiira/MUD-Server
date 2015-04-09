package main

import (
	"encoding/xml"
	"github.com/daviddengcn/go-colortext"
	"strconv"
)

type Inventory struct {
	itemsI map[string]Item_I
}

//=================== CONSTRUCTORS =====================//

func newInvXML() *InventoryXML {
	invXML := new(InventoryXML)
	invXML.Items = make([]ItemXML, 10)
	invXML.Weapons = make([]WeaponXML, 10)
	invXML.Armours = make([]ArmourXML, 10)

	return invXML
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

func newInventory() *Inventory {
	i := new(Inventory)
	i.itemsI = make(map[string]Item_I)

	return i
}

//================== CLASS FUNCTIONS =============//

func (inv *Inventory) addItemToInventory(item Item_I) {
	if val, ok := inv.itemsI[item.getName()]; ok { // the item is already there
		val.increaseQuantity()
	} else {
		inv.itemsI[item.getName()] = item
	}
}
func (inv *Inventory) addWeaponToInventory(weapon *Weapon) {
	if val, ok := inv.itemsI[weapon.Item.name]; ok { // the item is already there
		val.increaseQuantity()
	} else {
		inv.itemsI[weapon.name] = weapon
	}
}
func (inv *Inventory) addArmourToInventory(armr *Armour) {
	if val, ok := inv.itemsI[armr.name]; ok { // the item is already there
		val.increaseQuantity()
	} else {
		inv.itemsI[armr.name] = armr
	}
}

func (inv *Inventory) getItemByName(name string) (Item_I, bool) {
	itm, found := inv.itemsI[name]
	return itm, found
}

func (inv *Inventory) getInventoryDescription() []FormattedString {
	//TODO getInventoryDescription
	return nil
}

func (inv *Inventory) toXML() *InventoryXML {
	invXML := newInvXML()

	for _, item := range inv.itemsI {
		invXML.Items = append(invXML.Items, item.toXML())
	}

	return invXML
}

type InventoryXML struct {
	XMLName xml.Name    `xml:"Inventory"`
	ItemsI  []ItemXML_I `xml:"ItemsI"`
	Items   []ItemXML   `xml:"Item"`
	Weapons []WeaponXML `xml:"Weapon"`
	Armours []ArmourXML `xml:"Armour"`
}
