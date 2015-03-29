package main

type Item struct {
	name        string
	description string
	itemLevel   int
	itemWorth   int

	//	itemID      int
	// properties
}

func (i *Item) getName() string {
	return i.name
}

func (i *Item) getDescription() string {
	return i.description
}
