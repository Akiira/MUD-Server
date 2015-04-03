package main

type Agent struct {
	Name         string
	currentHP    int
	MaxHitPoints int
	RoomIN       int
	//	PersonalInv    Inventory
	//	EquippedArmour ArmourSet

	//Core Stats
	Strength     int
	Constitution int
	Dexterity    int
	Wisdom       int
	Charisma     int
	Inteligence  int
}

//The go community says to end interface names with "er"
type Agenter interface {
	makeAttack(target Agenter) []FormattedString
	takeDamage(amount int, typeOfDamge int) []FormattedString
	getName() string
	getDefense() int
	isDead() bool
}
