package main

import (
	"encoding/xml"
)

type Weapon struct {
	Item
	attack int
	damage int
}

type WeaponXML struct {
	XMLName  xml.Name `xml:"Weapon"`
	ItemInfo ItemXML  `xml:"Item"`
	Attack   int      `xml:"Attack"`
	Damage   int      `xml:"Damage"`
}
