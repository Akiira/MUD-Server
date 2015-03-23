package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct {
	playerORMonster bool
	nameOfAgent     string
	action          string
	valueOrTarget   string
}

func newEvent(agent bool, nameOfAgent string, action string, value string) Event {
	event := new(Event)
	event.playerORMonster = agent
	event.nameOfAgent = nameOfAgent
	event.action = action
	event.valueOrTarget = value
	return *event
}
