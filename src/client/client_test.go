package redis

import (
	"io"
	"redis-go/server"
	"testing"
)

func TestClient_Close(t *testing.T) {
	// Client.Close() should properly close the client connection
	serverConn, client := MockServerClient()
	err := client.Close()
	if err != nil {
		t.Fatalf("error during connection close - %s", err)
	}
	if _, err = serverConn.Read(make([]byte, 1)); err != io.EOF { // error if connection has been closed
		t.Fatalf("connection was not closed properly: expect io.EOF error, got %s", err)
	}
}

func TestClient_Ping(t *testing.T) {
	serverConn, client := MockServerClient()
	go server.NewServer().Serve(serverConn)
	response, err := client.Ping()
	if err != nil {
		t.Fatalf("error during call to Client.Ping - %s", err)
	}
	if response != "PONG" {
		t.Fatalf("expected 'PONG' response, got %s", response)
	}
}

func TestClient_Echo(t *testing.T) {
	serverConn, client := MockServerClient()
	go server.NewServer().Serve(serverConn)
	response, err := client.Echo("toto")
	if err != nil {
		t.Fatalf("error during call to Client.Ping - %s", err)
	}
	if response != "toto" {
		t.Fatalf("expected 'toto' response, got %s", response)
	}
}
