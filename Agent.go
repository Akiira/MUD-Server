package main

type Agent interface {
	name string
	health int
	personalInv Inventory
	equippedArmour ArmourSet
	
	makeAttack(target *Agent) []FormattedString
	takeDamage(amount int, typeOfDamge int) []FormattedString
	getDefense() int
	getAttack() int
}
