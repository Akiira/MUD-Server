package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math/rand"
	"strconv"
)

type Weapon struct {
	Item
	attack int
	minDmg int
	maxDmg int
}

func (wpn *Weapon) GetType() int {
	return WEAPON
}

func (wpn *Weapon) GetAttack() int {
	return wpn.attack
}

func (wpn *Weapon) GetDamage() int {
	return rand.Intn(wpn.GetDamageRange()) + wpn.minDmg
}

func (wpn *Weapon) GetDamageRange() int {
	return wpn.maxDmg - wpn.minDmg + 1
}

func (w *Weapon) GetCopy() Item_I {
	wpn := new(Weapon)
	*wpn = *w
	return wpn
}

func (w *Weapon) GetWeaponPage() []FormattedString {
	output := newFormattedStringCollection()
	output.addMessage(ct.Green, "\t\t\tEquipped Weapon\n")
	output.addMessage2(fmt.Sprintf("\t%-15s   %-15s %-15s %-15s\n", "Name", "Attack", "MinDmg", "MaxDmg"))
	output.addMessage(ct.Green, "--------------------------------------------------------------------\n")
	output.addMessage2(fmt.Sprintf("\t%-15s   %-15s %-15s %-15s\n\n", w.name, strconv.Itoa(w.attack), strconv.Itoa(w.minDmg), strconv.Itoa(w.maxDmg)))

	return output.fmtedStrings
}

type WeaponXML struct {
	XMLName  xml.Name `xml:"Weapon"`
	ItemInfo *ItemXML `xml:"Item"`
	Attack   int      `xml:"Attack"`
	MinDmg   int      `xml:"MinDmg"`
	MaxDmg   int      `xml:"MaxDmg"`
}

func NewWeaponFromXML(weaponData *WeaponXML) *Weapon {
	wpn := new(Weapon)
	wpn.Item = *NewItemFromXML(weaponData.ItemInfo)
	wpn.attack = weaponData.Attack
	wpn.minDmg = weaponData.MinDmg
	wpn.maxDmg = weaponData.MaxDmg

	return wpn
}

func (w *Weapon) ToXML() ItemXML_I {
	wpnXML := new(WeaponXML)
	wpnXML.ItemInfo = w.Item.ToXML().(*ItemXML)
	wpnXML.Attack = w.attack
	wpnXML.MinDmg = w.minDmg
	wpnXML.MaxDmg = w.maxDmg

	return wpnXML
}

func (w WeaponXML) ToItem() Item_I {
	return NewWeaponFromXML(&w)
}
