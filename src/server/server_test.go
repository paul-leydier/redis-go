package server

import (
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

func Test_ConcurrentClients(t *testing.T) {
	// Multiple clients should be able to interact concurrently with the server
	go Run("tcp", "localhost", "6380")
	client1 := redis.NewClient("localhost", "6380")
	client2 := redis.NewClient("localhost", "6380")
	_, err := client1.Ping()
	if err != nil {
		t.Fatalf("error during client1 Ping - %s", err)
	}
	_, err = client2.Ping()
	if err != nil {
		t.Fatalf("error during client2 Ping - %s", err)
	}
	err = client1.Close()
	if err != nil {
		t.Fatalf("error while closing client1 - %s", err)
	}
	err = client2.Close()
	if err != nil {
		t.Fatalf("error while closing client1 - %s", err)
	}
}

func Test_Ping(t *testing.T) {
	// A "PING" command should receive a "PONG" response
	serverConn, client := redis.MockServerClient()
	go func() {
		Serve(serverConn)
	}()
	resp, err := client.Ping()
	if err != nil {
		t.Fatalf("error during client.Ping - %s", err)
	}
	if resp != "PONG" {
		t.Fatalf("invalid response: expected %s, got %s", "PONG", resp)
	}
}

func Test_Multiple_Pings(t *testing.T) {
	// A single connection should be able to send multiple commands
	serverConn, client := redis.MockServerClient()
	go func() {
		Serve(serverConn)
	}()
	for i := 0; i < 10; i++ {
		resp, err := client.Ping()
		if err != nil {
			t.Fatalf("error during client.Ping - %s", err)
		}
		if resp != "PONG" {
			t.Fatalf("invalid response: expected %s, got %s", "PONG", resp)
		}
	}
}

func TestInvalidCommand(t *testing.T) {
	serverConn, client := redis.MockServerClient()
	go Serve(serverConn)
	err := client.CustomCommand("FOO")
	if err != nil {
		t.Fatalf("error while sending invalid command - %s", err)
	}
	_, err = client.SimpleStringResponse()
	if err == nil {
		t.Fatalf("invalid command did not raise an error")
	}
}
