package main

import (
	"bytes"
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

func Test_Ping(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		serve(server)
	}()
	msg := []byte("PING")
	_, err := client.Write(msg)
	if err != nil {
		t.Fatalf("could not write to localhost:6379 - %s", err)
	}
	resp := make([]byte, 64)
	_, err = client.Read(resp)
	resp = bytes.Trim(resp, "\x00")
	if err != nil {
		t.Fatalf("could not read from localhost:6379 - %s", err)
	}
	if string(resp) != "+PONG\r\n" {
		t.Fatalf("invalid response: expected %s, got %s", "+PONG", resp)
	}
}
