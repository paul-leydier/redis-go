package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	network string = "tcp"
	url     string = "localhost"
	port    string = "6379"
)

func main() {
	l, err := net.Listen(network, url+":"+port)
	if err != nil {
		log.Fatalf("could not bind to port - %s", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("could not accept connection - %s", err)
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	defer conn.Close()
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
	command := strings.Split(strings.ToUpper(strings.TrimSpace(string(msg))), " ")
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
		return respEncode("string", command[1])
	}
	return respEncode("string", "PONG")
}

func respEncode(responseType string, content string) []byte {
	switch responseType {
	case "string":
		return []byte(fmt.Sprintf("+%s\r\n", content))
	default:
		return []byte("")
	}
}
