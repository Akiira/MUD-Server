// Armour
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"strconv"
	"strings"
)

type Armour struct {
	Item
	defense      int
	wearLocation string
}

//------------------- CONSTRUCTORS ------------------------//

func NewArmour(name1 string, descr string, def int, wearLoc string) Armour {
	a := Armour{defense: def, wearLocation: strings.ToLower(wearLoc)}
	a.name = name1
	a.description = descr

	return a
}

func NewArmourFromXML(armourData *ArmourXML) *Armour {
	arm := new(Armour)
	arm.Item = *NewItem(armourData.ItemInfo)
	arm.defense = armourData.Defense
	arm.wearLocation = strings.ToLower(armourData.WearLocation)

	return arm
}

//------------------- GETTERS -----------------------------//

func (arm *Armour) GetType() int {
	return ARMOUR
}

func (arm *Armour) GetCopy() Item_I {
	armr := new(Armour)
	*armr = *arm

	return armr
}

func (arm *Armour) ToXML() ItemXML_I {
	armXML := new(ArmourXML)
	armXML.ItemInfo = arm.Item.ToXML().(*ItemXML)
	armXML.Defense = arm.defense
	armXML.WearLocation = arm.wearLocation

	return armXML
}

func (a ArmourXML) ToItem() Item_I {
	return NewArmourFromXML(&a)
}

//=================== ARMOURSET CLASS =====================//

var locations = [...]string{"head", "chest", "legs", "feet", "hands"}

type ArmourSet struct {
	equipedArmour map[string]*Armour
}

//------------------- CONSTRUCTORS ------------------------//

func NewArmourSet() *ArmourSet {
	as := new(ArmourSet)
	as.equipedArmour = make(map[string]*Armour, 5)
	return as
}

func NewArmourSetFromXML(armourSetData *ArmourSetXML) *ArmourSet {
	as := NewArmourSet()

	for _, arm := range armourSetData.ArmSet {
		as.EquipArmour(NewArmourFromXML(&arm))
	}

	return as
}

//=================== CLASS FUNCTIONS =====================//

func (as *ArmourSet) GetDefense() int {
	defense := 0

	for _, armr := range as.equipedArmour {
		defense += armr.defense
	}

	return defense
}

func (as *ArmourSet) GetAndRemoveArmour(nameOrLocation string) *Armour {

	if IsLocation(nameOrLocation) {
		loc := nameOrLocation

		arm := as.equipedArmour[loc]
		delete(as.equipedArmour, loc)

		return arm
	} else {
		name := nameOrLocation

		for loc, armr := range as.equipedArmour {
			lcName := strings.ToLower(armr.name)

			if strings.Contains(lcName, name) {
				as.GetAndRemoveArmour(loc)
				return armr
			}
		}
	}
	return nil
}

func (as *ArmourSet) EquipArmour(arm *Armour) {
	as.equipedArmour[arm.wearLocation] = arm
}

func (as *ArmourSet) GetArmourWornPage() []FormattedString {
	output := newFormattedStringCollection()
	output.addMessage(ct.Green, "\t\t\tEquipped Armour\n")
	output.addMessage2(fmt.Sprintf("\t%-20s   %-20s %-20s\n", "Location", "Name", "Defense"))
	output.addMessage(ct.Green, "--------------------------------------------------------------------\n")

	for _, loc := range locations {
		if arm, found := as.equipedArmour[loc]; found {
			output.addMessage2(fmt.Sprintf("\t%-20s   %-20s %-20s\n", loc, arm.name, strconv.Itoa(arm.defense)))
		} else {
			output.addMessage2(fmt.Sprintf("\t%-20s   %-20s %-20s\n", loc, " ", "0"))
		}
	}

	return output.fmtedStrings
}

func (as *ArmourSet) IsArmourAt(loc string) bool {
	_, present := as.equipedArmour[loc]
	return present
}

func (as *ArmourSet) ToXML() *ArmourSetXML {
	asXML := new(ArmourSetXML)

	for _, arm := range as.equipedArmour {
		asXML.ArmSet = append(asXML.ArmSet, *arm.ToXML().(*ArmourXML))
	}

	return asXML
}

func IsLocation(loc string) bool {
	for _, element := range locations {
		if loc == element {
			return true
		}
	}

	return false
}

//=================== XML STUFF =====================//

type ArmourSetXML struct {
	XMLName xml.Name    `xml:"ArmourSet"`
	ArmSet  []ArmourXML `xml:"Armour"`
}

type ArmourXML struct {
	XMLName      xml.Name `xml:"Armour"`
	ItemInfo     *ItemXML `xml:"Item"`
	Defense      int      `xml:"Defense"`
	WearLocation string   `xml:"Location"`
}
