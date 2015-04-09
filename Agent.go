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

func (a *Agent) setAgentStatsFromXML(charData *CharacterXML) {
	a.Strength = charData.Strength
	a.Wisdom = charData.Wisdom
	a.Inteligence = charData.Inteligence
	a.Dexterity = charData.Dexterity
	a.Charisma = charData.Charisma
	a.Constitution = charData.Constitution
	a.currentHP = charData.HP
	a.MaxHitPoints = charData.HP
	a.Name = charData.Name
	a.RoomIN = charData.RoomIN
}

//func (a *Agent) getDamage() {

//The go community says to end interface names with "er"
type Agenter interface {
	makeAttack(target string) []FormattedString
	takeDamage(amount int, typeOfDamge int) []FormattedString
	getName() string
	getDefense() int
	isDead() bool
	getRoomID() int
	//getClientConnection() *ClientConnection
}
