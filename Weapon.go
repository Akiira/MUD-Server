package main

import (
	"encoding/xml"
)

type Weapon struct {
	Item
	attack int
	damage int
}

func (wpn *Weapon) getAttack() int {
	return wpn.attack
}

func (wpn *Weapon) getDamage() int {
	return wpn.damage
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
	wpnXML := new(WeaponXML)
	wpnXML.ItemInfo = *w.Item.toXML()
	wpnXML.Attack = w.attack
	wpnXML.Damage = w.damage

	return wpnXML
}
