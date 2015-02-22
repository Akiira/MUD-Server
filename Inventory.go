package main

type Inventory struct {
	items         [100]Item
	numberOfItems int
}

func createInventory(items [100]Item) *Inventory {

	i := new(Inventory)
	i.items = items
	i.numberOfItems = 1

	return i
}

func (inv *Inventory) getItemByName(name string) Item {
	var i int
	for i = 0; i < inv.numberOfItems; i++ {
		if inv.items[i].name == name {
			return inv.items[i]
		}
	}
	var null Item
	return null
}
