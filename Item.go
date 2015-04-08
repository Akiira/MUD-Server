package main

import (
	"encoding/xml"
)

type Item struct {
	name        string
	description string
	itemLevel   int
	itemWorth   int
	quantity    int
}

type ItemXML struct {
	XMLName     xml.Name `xml:"Item"`
	Name        string   `xml:"Name"`
	Description string   `xml:"Description"`
	ItemLevel   int      `xml:"Level"`
	ItemWorth   int      `xml:"Worth"`
	Quantity    int      `xml:"Quantity"`
}

func itemFromXML(itemData *ItemXML) *Item {
	itm := new(Item)
	itm.description = itemData.Description
	itm.itemLevel = itemData.ItemLevel
	itm.itemWorth = itemData.ItemWorth
	itm.name = itemData.Name
	itm.quantity = itemData.Quantity

	return itm
}

func (i *Item) getName() string {
	return i.name
}

func (i *Item) getDescription() string {
	return i.description
}

func (i *Item) toXML() *ItemXML {
	xmlItem := new(ItemXML)
	xmlItem.Name = i.name
	xmlItem.Description = i.description
	xmlItem.ItemLevel = i.itemLevel
	xmlItem.ItemWorth = i.itemWorth
	xmlItem.Quantity = i.quantity

	return xmlItem
}
