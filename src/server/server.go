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
	respType, message := core.RespDecode(msg)
	if respType != core.SimpleString {
		return nil, fmt.Errorf("commands should be SimpleString, got %d - %s", respType, message)
	}
	command := strings.Split(strings.ToUpper(strings.TrimSpace(message)), " ")
	if len(command) == 0 {
		return nil, fmt.Errorf("empty command - %s", command)
	}
	switch command[0] {
	case "PING":
		return pingCommand(command), nil
	default:
		return nil, InvalidCommandError{received: command[0]}
	}
}

func pingCommand(command []string) []byte {
	if len(command) > 1 {
		return core.RespEncode(core.SimpleString, command[1])
	}
	return core.RespEncode(core.SimpleString, "PONG")
}

// Serving errors -----------------------------------------

func handleServingError(err error) []byte {
	if icerr, ok := err.(InvalidCommandError); ok {
		return core.RespEncode(core.Error, fmt.Sprintf("ERR %s", icerr.Error()))
	}
	log.Printf("serving error - %s", err)
	return core.RespEncode(core.Error, "ERR internal error")
}

type InvalidCommandError struct {
	received string
}

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf("unknown command '%s'", e.received)
}
