package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct {
	action string
	target string
	agent  Agenter
}

func NewEvent(agent Agenter, action string, target string) Event {
	event := new(Event)
	event.agent = agent
	event.action = action
	event.target = target
	return *event
}
