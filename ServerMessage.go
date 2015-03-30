// ServerMessage
package main

type ServerMessage struct {
	Value []FormattedString
}

func newServerMessage(msgs []FormattedString) ServerMessage {
	return ServerMessage{Value: msgs}
}

func getMessage(srvMsg ServerMessage) string {
	if len(srvMsg.Value) == 0 {
		return ""
	}
	return srvMsg.Value[0].Value
}

func (msg *ServerMessage) isError() bool {
	if len(msg.Value) == 0 {
		return false
	}

	return (strings.Split(message.Value, " ")[0] == "error")
}
