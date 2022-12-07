// Instructions logic -------------------------------------

package server

import (
	"fmt"
	"redis-go/core"
	"strings"
)

func handleMessage(msg []byte) ([]byte, error) {
	encoded := core.NewEncodedRespElem(msg)
	messages, err := encoded.Decode().Array()
	if err != nil {
		return nil, fmt.Errorf("received invalid message - %s", err)
	}
	command := make([]string, len(messages))
	for i, msg := range messages {
		command[i], err = msg.String()
		if err != nil {
			return nil, err
		}
	}
	if len(command) == 0 {
		return nil, fmt.Errorf("empty command - %s", command)
	}
	switch strings.ToUpper(command[0]) {
	case "PING":
		return pingCommand(command), nil
	case "ECHO":
		return echoCommand(command), nil
	default:
		return nil, InvalidCommandError{received: command[0]}
	}
}

func pingCommand(command []string) []byte {
	if len(command) > 1 {
		return core.RespElem{
			Type:    core.SimpleString,
			Content: command[1],
		}.Encode()
	}
	return core.RespElem{
		Type:    core.SimpleString,
		Content: "PONG",
	}.Encode()
}

func echoCommand(command []string) []byte {
	if len(command) == 1 {
		return []byte("")
	}
	return core.RespElem{
		Type:    core.SimpleString,
		Content: strings.Join(command[1:], " "),
	}.Encode()
}
