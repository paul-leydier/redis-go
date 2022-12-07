// Instructions logic -------------------------------------

package server

import (
	"errors"
	"fmt"
	"redis-go/core"
	"strings"
)

func (s *Server) handleMessage(msg []byte) ([]byte, error) {
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
	case "GET":
		value, err := s.getCommand(command)
		if err != nil {
			return nil, err
		}
		return value, nil
	case "SET":
		response, err := s.setCommand(command)
		if err != nil {
			return nil, err
		}
		return response, nil
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

func (s *Server) getCommand(command []string) ([]byte, error) {
	if len(command) < 2 {
		return nil, errors.New("need a key for the GET command")
	}
	return core.RespElem{
		Type:    core.BulkString,
		Content: s.storage.Get(command[1]),
	}.Encode(), nil
}

func (s *Server) setCommand(command []string) ([]byte, error) {
	if len(command) < 3 {
		return nil, fmt.Errorf("need a key and a value for the SET command - got %v", command)
	}
	s.storage.Set(command[1], command[2])
	return core.RespElem{
		Type:    core.SimpleString,
		Content: "OK",
	}.Encode(), nil
}
