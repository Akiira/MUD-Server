// Armour
package main

import (
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

//--------------ARMOURSET CLASS----------------

type ArmourSet struct {
	equipedArmour map[string]Armour
}

func newArmourSet() ArmourSet {
	as := new(ArmourSet)
	as.equipedArmour = make(map[string]Armour, 5)
	return *as
}

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
