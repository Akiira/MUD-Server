package main

import ()

type Trader struct {
	isSelected  bool
	isConfirmed bool
	itemList    []string
	quanList    []int
	dealer      *ClientConnection
}
