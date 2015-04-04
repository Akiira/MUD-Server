// Armour
package main

import (
	"encoding/xml"
	"github.com/daviddengcn/go-colortext"
)

type Armour struct {
	Item
	defense      int
	wearLocation string
}

//When should constructors return a pointer instead of the object itself?
func newArmour(name1 string, descr string, def int, wearLoc string) Armour {
	a := Armour{defense: def, wearLocation: wearLoc}
	a.name = name1
	a.description = descr
	return a
	//return &a
}

func armourFromXML(armourData ArmourXML) *Armour {
	arm := new(Armour)
	arm.Item = *itemFromXML(&armourData.ItemInfo)
	arm.defense = armourData.Defense
	arm.wearLocation = armourData.WearLocation

	return arm
}

func (arm *Armour) toXML() *ArmourXML {
	armXML := new(ArmourXML)
	armXML.ItemInfo = *arm.Item.toXML()
	armXML.Defense = arm.defense
	armXML.WearLocation = arm.wearLocation

	return armXML
}

//--------------ARMOURSET CLASS----------------

type ArmourSet struct {
	equipedArmour map[string]Armour
}

//=================== CONSTRUCTORS =====================//

func newArmourSet() ArmourSet {
	as := new(ArmourSet)
	as.equipedArmour = make(map[string]Armour, 5)
	return *as
}

func armourSetFromXML(armourSetData ArmourSetXML) *ArmourSet {
	as := new(ArmourSet)

	for _, arm := range armourSetData.ArmSet {
		as.equipArmour(arm.WearLocation, *armourFromXML(arm))
	}

	return as
}

//=================== CLASS FUNCTIONS =====================//

func (as *ArmourSet) getArmoursDefense() int {
	defense := 0

	for _, armr := range as.equipedArmour {
		defense += armr.defense
	}

	return defense
}

func (as *ArmourSet) takeOffArmourByLocation(loc string) Armour {
	//TODO add check for no armour worn at location
	arm := as.equipedArmour[loc]

	delete(as.equipedArmour, loc)

	return arm
}

func (as *ArmourSet) equipArmour(location string, arm Armour) {
	//TODO add check for armour already being worn
	as.equipedArmour[location] = arm
}

func (as *ArmourSet) getListOfArmourWorn() []FormattedString {
	foo := []string{"head", "chest", "legs", "feet", "hands"}
	var output []FormattedString
	for _, e := range foo {
		str := "\t" + e + ": "
		arm, present := as.equipedArmour[e]
		if present {
			str += arm.getName()
		}
		fmtStr := FormattedString{Color: ct.White, Value: str + "\n"}
		output = append(output, fmtStr)
	}

	return output
}

func (as *ArmourSet) isArmourEquippedAtLocation(loc string) bool {
	_, present := as.equipedArmour[loc]
	return present
}

func (as *ArmourSet) toXML() *ArmourSetXML {
	asXML := new(ArmourSetXML)

	for _, arm := range as.equipedArmour {
		asXML.ArmSet = append(asXML.ArmSet, *arm.toXML())
	}

	return asXML
}

//=================== XML STUFF =====================//

type ArmourSetXML struct {
	XMLName xml.Name    `xml:"ArmourSet"`
	ArmSet  []ArmourXML `xml:"Armour"`
}

type ArmourXML struct {
	XMLName      xml.Name `xml:"Armour"`
	ItemInfo     ItemXML  `xml:"Item"`
	Defense      int      `xml:"Defense"`
	WearLocation string   `xml:"Location"`
}
