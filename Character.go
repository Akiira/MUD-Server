package main

type Character struct {
	Name string
	RoomIN int
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
