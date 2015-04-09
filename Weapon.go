package main

import (
	"encoding/xml"
	"math/rand"
)

type Weapon struct {
	Item
	attack int
	minDmg int
	maxDmg int
}

func (wpn *Weapon) getItemType() int {
	return WEAPON
}

func (wpn *Weapon) getAttack() int {
	return wpn.attack
}

func (wpn *Weapon) getDamage() int {
	return rand.Intn(wpn.getDamageRange()) + wpn.minDmg
}

func (wpn *Weapon) getDamageRange() int {
	return wpn.maxDmg - wpn.minDmg + 1
}

type WeaponXML struct {
	XMLName  xml.Name `xml:"Weapon"`
	ItemInfo ItemXML  `xml:"Item"`
	Attack   int      `xml:"Attack"`
	MinDmg   int      `xml:"MinDmg"`
	MaxDmg   int      `xml:"MaxDmg"`
}

func weaponFromXML(weaponData *WeaponXML) *Weapon {
	wpn := new(Weapon)
	wpn.Item = *itemFromXML(&(weaponData.ItemInfo))
	wpn.attack = weaponData.Attack
	wpn.minDmg = weaponData.MinDmg
	wpn.maxDmg = weaponData.MaxDmg

	return wpn
}

func (w *Weapon) toXML() *WeaponXML {
	wpnXML := new(WeaponXML)
	wpnXML.ItemInfo = *w.Item.toXML()
	wpnXML.Attack = w.attack
	wpnXML.MinDmg = w.minDmg
	wpnXML.MaxDmg = w.maxDmg

	return wpnXML
}
