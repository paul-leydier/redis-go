package main

import "redis-go/server"

const (
	network string = "tcp"
	url     string = "localhost"
	port    string = "6379"
)

func main() {
	s := server.NewServer()
	s.Run(network, url, port)
}
