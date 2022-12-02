package redis

import (
	"fmt"
	"net"
	"redis-go/core"
)

const network string = "tcp"

type serverInfo struct {
	Url  string
	Port string
}

type Client struct {
	server serverInfo
	conn   net.Conn
}

func NewClient(url string, port string) Client {
	client := Client{
		server: serverInfo{
			Url:  url,
			Port: port,
		},
	}
	return client
}

func (r *Client) Connect() error {
	conn, err := net.Dial(network, r.server.Url+":"+r.server.Port)
	r.conn = conn
	return err
}

func (r *Client) Close() error {
	return r.conn.Close()
}

func (r *Client) CustomCommand(command string, content string) error {
	_, err := r.conn.Write(core.RespEncode(core.SimpleString, command+" "+content))
	return err
}

func (r *Client) SimpleStringResponse() (string, error) {
	encodedResponse := make([]byte, 64)
	_, err := r.conn.Read(encodedResponse)
	if err != nil {
		return "", err
	}
	respType, resp := core.RespDecode(encodedResponse)
	if respType != core.SimpleString {
		return "", fmt.Errorf("invalid response - %s", resp)
	}
	return resp, nil
}

func (r *Client) Ping(content string) (string, error) {
	err := r.CustomCommand("PING", content)
	if err != nil {
		return "", err
	}
	return r.SimpleStringResponse()
}

func MockServerClient() (net.Conn, Client) {
	clientConn, serverConn := net.Pipe()
	client := NewClient("localhost", "6379")
	client.conn = clientConn
	return serverConn, client
}
