package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io"
)

type entry struct {
	item     Item_I
	quantity int
}
type Inventory struct {
	items map[string]*entry
}

//=================== CONSTRUCTORS =====================//

func newInvXML() *InventoryXML {
	invXML := new(InventoryXML)
	invXML.Items = make([]interface{}, 10)

	return invXML
}

func inventoryFromXML(invXml *InventoryXML) *Inventory {
	inv := newInventory()

	for _, itm := range invXml.Items {
		inv.addItemToInventory(itm.(ItemXML_I).toItem())
	}

	return inv
}

func newInventory() *Inventory {
	i := new(Inventory)
	i.items = make(map[string]*entry)

	return i
}

//================== CLASS FUNCTIONS =============//

func (inv *Inventory) addItemToInventory(item Item_I) {
	if val, ok := inv.items[item.getName()]; ok { // the item is already there
		val.quantity++
	} else {
		inv.items[item.getName()] = &entry{item: item, quantity: 1}
	}
}

//PossesItem returns true if the inventory contains an item with the given name, else false.
func (inv *Inventory) PossesItem(name string) bool {
	_, found := inv.GetItem(name)
	return found
}

//GetAndRemoveItem returns a pointer to the item in the inventory with the specified
//name and removes it from the inventory.
//If the item is not found the pointer is nil and bool is false.
func (inv *Inventory) GetAndRemoveItem(name string) (Item_I, bool) {

	item, found := inv.items[name]
	if found {
		item.quantity--
		if item.quantity == 0 {
			delete(inv.items, name)
		}

		return item.item, true
	} else {
		return nil, false
	}
}

//GetItem returns a pointer to the item in the inventory with the name in name.
//If the item is found it is not removed from the inventory.
//If the item is not found the pointer is nil and bool is false.
func (inv *Inventory) GetItem(name string) (Item_I, bool) {
	item, found := inv.items[name]

	return item.item, found
}

func (inv *Inventory) getInventoryDescription() []FormattedString {
	desc := newFormattedStringCollection()
	desc.addMessage2("\nInventory\n")
	desc.addMessage(ct.Green, "-----------------------------------------\n")

	for name, itemEntry := range inv.items {
		desc.addMessage(ct.Green, fmt.Sprintf("\t%-20s   %3d", name, itemEntry.quantity)+"\n")
	}
	desc.addMessage2("\n")
	return desc.fmtedStrings
}

func (inv *Inventory) toXML() *InventoryXML {
	invXML := newInvXML()

	for _, item := range inv.items {
		for i := 0; i < item.quantity; i++ {
			invXML.Items = append(invXML.Items, item.item.toXML())
		}
	}

	return invXML
}

//================== XML STUFF =============//

type InventoryXML struct {
	XMLName xml.Name      `xml:"Inventory"`
	Items   []interface{} `xml:",any"`
}

func (c *InventoryXML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var item ItemXML_I

	for t, err := d.Token(); err != io.EOF; {
		switch t1 := t.(type) {
		case xml.StartElement:
			if t1.Name.Local == "Armour" {
				item = new(ArmourXML)
			} else if t1.Name.Local == "Weapon" {
				item = new(WeaponXML)
			} else {
				item = new(ItemXML)
			}

			err = d.DecodeElement(item, &t1)
			checkError(err, true)
			c.Items = append(c.Items, item)
		}

		t, err = d.Token()
	}

	return nil
}
