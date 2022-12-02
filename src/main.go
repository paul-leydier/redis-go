package main

import "redis-go/server"

const (
	network string = "tcp"
	url     string = "localhost"
	port    string = "6379"
)

func main() {
	server.Run(network, url, port)
}
