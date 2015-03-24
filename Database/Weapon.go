package main

type Weapon struct {
	Item
	attack int
}

func newWeapon(name1 string, descr string, atk int) Weapon {
	a := Weapon{attack: atk}
	a.name = name1
	a.description = descr
	return a
	//return &a
}
