package server

import (
	"bytes"
	"net"
	redis "redis-go/client"
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
	msg := []byte("+PING\r\n")
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
	msg := []byte("+PING\r\n")
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

func TestInvalidCommand(t *testing.T) {
	serverConn, client := redis.MockServerClient()
	go Serve(serverConn)
	err := client.CustomCommand("FOO", "")
	if err != nil {
		t.Fatalf("error while sending invalid command - %s", err)
	}
	_, err = client.SimpleStringResponse()
	if err == nil {
		t.Fatalf("invalid command did not raise an error")
	}
}
