package main

type Item struct {
	name        string
	description string
	itemID      int
}

func (i *Item) getName() string {
	return i.name
}

func (i *Item) getDescription() string {
	return i.description
}

func (i *Item) getItemID() int {
	return i.itemID
}
