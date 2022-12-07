package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"redis-go/core"
	"strings"
)

// Core server logic --------------------------------------

func Run(network string, url string, port string) {
	l, err := net.Listen(network, url+":"+port)
	if err != nil {
		log.Fatalf("could not bind to port - %s", err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Panicf("error while closing the net.Listener - %s", err)
		}
	}(l)
	Listen(l)
}
func Listen(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("could not accept connection - %s", err)
		}
		go Serve(conn)
	}
}

func Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("could not close connection - %s", err)
		}
	}(conn)
	for {
		msg := make([]byte, 64)
		_, err := conn.Read(msg)
		if err != nil {
			log.Fatalf("could not read message - %s", err)
		}
		msg = bytes.Trim(msg, "\x00")
		response, err := handleMessage(msg)
		if err != nil {
			response = handleServingError(err)
		}
		_, err = conn.Write(response)
		if err != nil {
			log.Fatalf("could not write response - %s", err)
		}
	}
}

// Instructions logic -------------------------------------

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
		Content: command[1],
	}.Encode()
}

// Serving errors -----------------------------------------

func handleServingError(err error) []byte {
	if icerr, ok := err.(InvalidCommandError); ok {
		return core.RespElem{
			Type:    core.Error,
			Content: fmt.Sprintf("ERR %s", icerr.Error()),
		}.Encode()
	}
	log.Printf("serving error - %s", err)
	return core.RespElem{
		Type:    core.Error,
		Content: "ERR internal error",
	}.Encode()
}

type InvalidCommandError struct {
	received string
}

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf("unknown command '%s'", e.received)
}
