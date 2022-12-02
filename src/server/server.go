package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"redis-go/core"
	"strings"
)

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
			log.Fatalf("could not handle message - %s", err)
		}
		_, err = conn.Write(response)
		if err != nil {
			log.Fatalf("could not write response - %s", err)
		}
	}
}

func handleMessage(msg []byte) ([]byte, error) {
	command := strings.Split(strings.ToUpper(strings.TrimSpace(core.RespDecode(msg))), " ")
	if len(command) == 0 {
		return nil, fmt.Errorf("empty command - %s", command)
	}
	switch command[0] {
	case "PING":
		return pingCommand(command), nil
	default:
		return nil, fmt.Errorf("unknown command - %s", command)
	}
}

func pingCommand(command []string) []byte {
	if len(command) > 1 {
		return core.RespEncode(core.SimpleString, command[1])
	}
	return core.RespEncode(core.SimpleString, "PONG")
}
