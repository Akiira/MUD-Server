package main

type Character struct {
	name string
}

func (c *Character) getName() string {
	return c.name
}
