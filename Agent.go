package main

type Agenter interface {
	Attack(target Agenter) []FormattedString
	TakeDamage(amount int, typeOfDamge int)
	Respawn() []FormattedString
	AddTarget(target Agenter)

	GetName() string
	GetDescription() string
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

func (a *Agent) SetAgentStats(charData *CharacterXML) {
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
