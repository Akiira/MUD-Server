package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io"
)

type Inventory struct {
	items map[string][]Item_I
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
		inv.AddItem(itm.(ItemXML_I).toItem())
	}

	return inv
}

func newInventory() *Inventory {
	i := new(Inventory)
	i.items = make(map[string][]Item_I)

	return i
}

//================== CLASS FUNCTIONS =============//

//AddInventory will ad all items from the given inventory into this one.
func (inv *Inventory) AddInventory(otherInv Inventory) {
	for _, items := range otherInv.items {
		inv.AddItems(items)
	}
}

//AddItems will add all items in the provided splice to the inventory.
func (inv *Inventory) AddItems(items []Item_I) {
	for _, item := range items {
		inv.AddItem(item)
	}
}

//AddItem will add the supplied item to the inventory.
func (inv *Inventory) AddItem(item Item_I) {
	if val, ok := inv.items[item.getName()]; ok { // the item is already there
		val = append(val, item)
	} else {
		inv.items[item.getName()] = make([]Item_I, 0)
		inv.AddItem(item)
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
	item, found := inv.GetItem(name)

	if found {
		inv.RemoveItem(name)
		return item, true
	} else {
		return nil, false
	}
}

//RemoveItem removes an item with the specified name from the inventory.
//If the item was not present then no change is made.
func (inv *Inventory) RemoveItem(name string) {
	_, found := inv.GetItem(name)

	if found {
		items := inv.items[name]
		inv.items[name] = items[0 : len(items)-2]
	}
}

//GetItem returns a pointer to the item in the inventory with the name in name.
//If the item is found it is not removed from the inventory.
//If the item is not found the pointer is nil and bool is false.
func (inv *Inventory) GetItem(name string) (Item_I, bool) {
	items, _ := inv.items[name]

	if len(items) > 0 {
		return items[len(items)-1], true
	}
	return nil, false
}

func (inv *Inventory) getInventoryDescription() []FormattedString {
	desc := newFormattedStringCollection()
	desc.addMessage2("\nInventory\n")
	desc.addMessage(ct.Green, "-----------------------------------------\n")

	for name, itemEntry := range inv.items {
		desc.addMessage(ct.Green, fmt.Sprintf("\t%-20s   %3d", name, len(itemEntry))+"\n")
	}
	desc.addMessage2("\n")
	return desc.fmtedStrings
}

func (inv *Inventory) toXML() *InventoryXML {
	invXML := newInvXML()

	for _, items := range inv.items {
		for _, item := range items {
			invXML.Items = append(invXML.Items, item.toXML())
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
