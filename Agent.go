package main

type Agent struct {
	Name           string
	Health         int
	RoomIn         int
	PersonalInv    Inventory
	EquippedArmour ArmourSet
}

//The go community says to end interface names with "er"
type Agenter interface {
	makeAttack(targetName string) []FormattedString
	takeDamage(amount int, typeOfDamge int) []FormattedString
}
