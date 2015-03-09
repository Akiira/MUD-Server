package main

import (
	"math/rand"
)

// this should be a stub that hold a connection to a client
// works like a thread on its own
type Character struct {
	Name      string
	RoomIN    int
	HitPoints int
	Defense   int
	//	Strength int
	//	Constitution int
	//	Dexterity int
	//	Wisdom int
	//	Charisma int
	//	Inteligence int

	//	Location string

	//	Race string
	//	Class string

	//	PersonalInvetory Inventory

	//	Weapon Item
	//ArmourSet
}

func (c *Character) getAttackRoll() int {
	return rand.Int() % 6
}

func (c *Character) addItemToInventory(item Item) {

}

func (c *Character) equipItemFromGround(item Item) {

}

func (c *Character) equipItemFromInventory(itemName string) {

}

func (c *Character) moveCharacter(direction string) {

}

func (c *Character) getAttack() int {
	return -1
}

func (c *Character) getName() string {
	return c.Name
}
