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

func newEventFromMessage(msg ClientMessage, agent Agenter) Event {
	return newEvent(agent, msg.Command, msg.Value)
}

func newEvent(agent Agenter, action string, target string) Event {
	event := new(Event)
	event.agent = agent
	event.action = action
	event.target = target
	return *event
}
