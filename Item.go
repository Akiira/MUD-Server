package main

import (
	"encoding/xml"
)

type Item_I interface {
	GetName() string
	GetDescription() []FormattedString
	GetType() int
	GetCopy() Item_I
	ToXML() ItemXML_I
}

type ItemXML_I interface {
	ToItem() Item_I
}

const (
	BASE_ITEM = 0
	WEAPON    = 1
	ARMOUR    = 2
)

type Item struct {
	name        string
	description string
	itemLevel   int
	itemWorth   int
	typeOfItem  int
}

type ItemXML struct {
	XMLName     xml.Name `xml:"Item"`
	Name        string   `xml:"Name"`
	Description string   `xml:"Description"`
	ItemLevel   int      `xml:"Level"`
	ItemWorth   int      `xml:"Worth"`
}

func NewItemFromXML(itemData *ItemXML) *Item {
	itm := new(Item)
	itm.description = itemData.Description
	itm.itemLevel = itemData.ItemLevel
	itm.itemWorth = itemData.ItemWorth
	itm.name = itemData.Name

	return itm
}

func (i *Item) GetName() string {
	return i.name
}

func (i *Item) GetDescription() []FormattedString {
	return newFormattedStringSplice(i.description)
}

func (i *Item) GetType() int {
	return BASE_ITEM
}

func (i *Item) GetCopy() Item_I {
	itm := new(Item)
	*itm = *i
	return itm
}

func (i *Item) ToXML() ItemXML_I {
	xmlItem := new(ItemXML)
	xmlItem.Name = i.name
	xmlItem.Description = i.description
	xmlItem.ItemLevel = i.itemLevel
	xmlItem.ItemWorth = i.itemWorth

	return xmlItem
}

func (i ItemXML) ToItem() Item_I {
	return NewItemFromXML(&i)
}
