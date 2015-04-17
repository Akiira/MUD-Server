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
func (inv *Inventory) checkAvailableItem(itemIndex int, quantity int) (bool, string) {
	if itemIndex <= len(inv.items) {
		i := 1
		for key, entryItem := range inv.items {
			if i == itemIndex && entryItem.quantity >= quantity {
				return true, key
			}
			i++
		}
		return false, " "
	} else {
		return false, " "
	}

}

func (inv *Inventory) checkAvailableItemMap(itemMap map[string]int) bool {

	for key, itemQuan := range itemMap {
		item, found := inv.items[key]
		if !found {
			return false
		} else if item.quantity < itemQuan {
			return false
		}
	}
	return true

}

func (inv *Inventory) removeItemMapFromInventory(itemMap map[string]int) []entry {

	var deductedItem []entry

	for key, itemQuan := range itemMap {

		tradeItem, _ := inv.items[key]

		deductedItem = append(deductedItem, entry{item: tradeItem.item.getCopy(), quantity: itemQuan})

		if tradeItem.quantity > itemQuan {
			tradeItem.quantity -= itemQuan
			inv.items[key] = tradeItem
		} else {
			delete(inv.items, key)
		}

	}
	return deductedItem
}

func (inv *Inventory) addItemListToInventory(entryList []entry) {
	for _, itemEntry := range entryList {
		tradeItem, found := inv.items[itemEntry.item.getName()]
		if found {
			tradeItem.quantity += itemEntry.quantity
			inv.items[itemEntry.item.getName()] = tradeItem
		} else {
			inv.items[itemEntry.item.getName()] = &itemEntry
		}
	}
}

func (inv *Inventory) addItemToInventory(item Item_I) {
	if val, ok := inv.items[item.getName()]; ok { // the item is already there
		val.quantity++
	} else {
		inv.items[item.getName()] = &entry{item: item, quantity: 1}
	}
}

func (inv *Inventory) PossesItem(name string) bool {
	_, found := inv.getItemByName(name)
	return found
}

func (inv *Inventory) getItemByName(name string) (Item_I, bool) {

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

func (inv *Inventory) getInventoryDescription() []FormattedString {
	desc := newFormattedStringCollection()
	desc.addMessage2("\nInventory\n")
	desc.addMessage(ct.Green, "-----------------------------------------\n")

	i := 1
	for name, itemEntry := range inv.items {
		desc.addMessage(ct.Green, fmt.Sprintf("%d\t%-20s   %3d", i, name, itemEntry.quantity)+"\n")
		i++
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
