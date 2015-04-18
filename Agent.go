package main

type Agenter interface {
	makeAttack(target Agenter) []FormattedString
	takeDamage(amount int, typeOfDamge int)
	respawn() *FmtStrCollection
	addTarget(target Agenter)

	GetName() string
	GetDefense() int
	GetRoomID() int

	IsDead() bool

	SendMessage(interface{})
}

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
