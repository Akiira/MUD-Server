package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct {
	action string
	target Agenter //TODO consider changing this to an Ageneter as well
	agent  Agenter
	client *ClientConnection
}

func newEventFromMessage(msg ClientMessage, agent Agenter, cc *ClientConnection) Event {
	target := eventManager.worldRooms[cc.getCharactersRoomID()].getMonster(msg.Value)
	return newEvent(agent, msg.Command, target, cc)
}

func newEvent(agent Agenter, action string, target Agenter, cc *ClientConnection) Event {
	event := new(Event)
	event.agent = agent
	event.action = action
	event.target = target
	event.client = cc
	return *event
}
