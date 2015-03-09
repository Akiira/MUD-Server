// EventManager
package main

import (
	"github.com/daviddengcn/go-colortext"
)

//event manager should only receive event from either monster / player and echo to all that monster / player in the room
// then those player / monster will decide by themselve to get hit or not
// with this concept of oop it should let us handle both eventmanager and play easily
type EventManager struct {
	//a Q of events
	//a timer
}
type FormattedString struct {
	Color ct.Color
	Value string
}

