package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct {
	action string
	target string
	agent  Agenter
	client *ClientConnection
}

func newEventFromMessage(msg ClientMessage, agent Agenter, cc *ClientConnection) Event {
	return newEvent(agent, msg.Command, msg.Value, cc)
}

func newEvent(agent Agenter, action string, target string, cc *ClientConnection) Event {
	event := new(Event)
	event.agent = agent
	event.action = action
	event.target = target
	event.client = cc
	return *event
}
