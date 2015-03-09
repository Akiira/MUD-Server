// Armour
package main

type Armour struct {
	Item
	defense int
	wearLocation int
}

//When should constructors return a pointer instead of the object itself?
func newArmour(name1 string, descr string, def int, wearLoc int) Armour {
	a := Armour{defense: def, wearLocation: wearLoc}
	a.name = name1
	a.description = descr
	return a
	//return &a
}