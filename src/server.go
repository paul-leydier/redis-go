package main

import (
	"log"
	"net"
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
	_, err = l.Accept()
	if err != nil {
		log.Fatalf("could not accept connetion - %s", err)
	}
}
