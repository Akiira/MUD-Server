package main

type Agent struct {
	name           string
	health         int
	RoomIn         int
	personalInv    Inventory
	equippedArmour ArmourSet
}

//The go community says to end interface names with "er"
type Agenter interface {
	makeAttack(targetName string) []FormattedString
	takeDamage(amount int, typeOfDamge int) []FormattedString
	getDefense() int
	getAttackRoll() int
}
