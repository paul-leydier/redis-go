package main

import (
	"net"
	"testing"
)

func Test_Listen(t *testing.T) {
	// should be able to bind to localhost:6379
	go main()
	_, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("could not connect to localhost:6379 - %s", err)
	}
}
