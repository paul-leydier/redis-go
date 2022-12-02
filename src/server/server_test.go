package server

import (
	"bytes"
	"net"
	"testing"
)

func Test_Listen(t *testing.T) {
	// should be able to bind to localhost:6379
	go Run("tcp", "localhost", "6379")
	_, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("could not connect to localhost:6379 - %s", err)
	}
}

func Test_Ping(t *testing.T) {
	// A "PING" command should receive a "PONG" response
	server, client := net.Pipe()
	go func() {
		Serve(server)
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

func Test_Multiple_Pings(t *testing.T) {
	// A single connection should be able to send multiple commands
	server, client := net.Pipe()
	go func() {
		Serve(server)
	}()
	msg := []byte("PING")
	for i := 0; i < 10; i++ {
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
}
