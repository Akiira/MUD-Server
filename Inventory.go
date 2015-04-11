package main

import (
	"encoding/xml"
	//"github.com/daviddengcn/go-colortext"
)

type Inventory struct {
	items map[string]Item_I
}

//=================== CONSTRUCTORS =====================//

func newInvXML() *InventoryXML {
	invXML := new(InventoryXML)
	invXML.Items = make([]ItemXML_I, 10)

	return invXML
}

func inventoryFromXML(invXml *InventoryXML) *Inventory {
	inv := newInventory()

	//loop through items
	for _, itm := range invXml.Items {
		inv.addItemToInventory(itm.toItem())
	}

	return inv
}

func newInventory() *Inventory {
	i := new(Inventory)
	i.items = make(map[string]Item_I)

	return i
}

//================== CLASS FUNCTIONS =============//

func (inv *Inventory) addItemToInventory(item Item_I) {
	if val, ok := inv.items[item.getName()]; ok { // the item is already there
		val.increaseQuantity()
	} else {
		inv.items[item.getName()] = item
	}
}

func (inv *Inventory) getItemByName(name string) (Item_I, bool) {
	itm, found := inv.items[name]
	return itm, found
}

func (inv *Inventory) getInventoryDescription() []FormattedString {
	//TODO getInventoryDescription
	return nil
}

func (inv *Inventory) toXML() *InventoryXML {
	invXML := newInvXML()

	for _, item := range inv.items {
		invXML.Items = append(invXML.Items, item.toXML())
	}

	return invXML
}

type InventoryXML struct {
	XMLName xml.Name    `xml:"Inventory"`
	Items   []ItemXML_I `xml:"Items"`
}
