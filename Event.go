package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct { //TODO use an Agent interface to avoid two pointers
	playerORMonster bool
	action          string
	valueOrTarget   string //TODO consider changing this to an Ageneter as well
	agent           Agenter
	client          *ClientConnection
}

//TODO fix these functions to use only *Agent rahter then charName, *Character or *Monster
//func newEventFromMessage(msg ClientMessage, charName string) Event {
//	return newEvent(PLAYER, charName, msg.Command, msg.Value)
//}
func newEventFromMessage(msg ClientMessage, agent Agenter, cc *ClientConnection) Event {
	return newEvent(PLAYER, agent, msg.Command, msg.Value, cc)
}

func newEvent(playerORMonster bool, agent Agenter, action string, value string, cc *ClientConnection) Event {
	event := new(Event)
	event.playerORMonster = playerORMonster
	event.agent = agent
	event.action = action
	event.valueOrTarget = value
	event.client = cc
	return *event
}
