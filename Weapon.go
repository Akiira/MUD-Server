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

func weaponFromXML(weaponData *WeaponXML) *Weapon {
	wpn := new(Weapon)
	wpn.Item = *itemFromXML(&weaponData.ItemInfo)
	wpn.attack = weaponData.Attack
	wpn.damage = weaponData.Damage

	return wpn
}

func (w *Weapon) toXML() *WeaponXML {
	//TODO toXML
	return nil
}
