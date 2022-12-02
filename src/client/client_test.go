package redis

import (
	"io"
	"net"
	"testing"
)

func mockServerClient() (net.Conn, Client) {
	clientConn, serverConn := net.Pipe()
	client := NewClient("localhost", "6379")
	client.conn = clientConn
	return serverConn, client
}

func TestClient_Close(t *testing.T) {
	// Client.Close() should properly close the client connection
	server, client := mockServerClient()
	err := client.Close()
	if err != nil {
		t.Fatalf("error during connection close - %s", err)
	}
	if _, err = server.Read(make([]byte, 1)); err != io.EOF { // error if connection has been closed
		t.Fatalf("connection was not closed properly: expect io.EOF error, got %s", err)
	}
}
