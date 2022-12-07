package redis

import (
	"fmt"
	"net"
	"redis-go/core"
	"strings"
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

func (r *Client) CustomCommand(command string) error {
	if r.conn == nil {
		err := r.Connect()
		if err != nil {
			return fmt.Errorf("could not connect to the Redis server - %s", err)
		}
	}
	words := strings.Split(command, " ")
	cmds := make([]core.RespElem, len(words))
	for i, word := range words {
		cmds[i] = core.RespElem{
			Type:    core.BulkString,
			Content: word,
		}
	}
	cmd := core.RespElem{
		Type:    core.Array,
		Content: cmds,
	}
	_, err := r.conn.Write(cmd.Encode())
	return err
}

func (r *Client) SimpleStringResponse() (string, error) {
	msg := make([]byte, 64)
	_, err := r.conn.Read(msg)
	if err != nil {
		return "", err
	}
	encoded := core.NewEncodedRespElem(msg)
	resp, err := encoded.Decode().String()
	if err != nil {
		return "", fmt.Errorf("invalid response - %s", err)
	}
	return resp, nil
}

func (r *Client) Ping() (string, error) {
	err := r.CustomCommand("PING")
	if err != nil {
		return "", err
	}
	return r.SimpleStringResponse()
}

func (r *Client) Echo(content string) (string, error) {
	err := r.CustomCommand("ECHO" + " " + content)
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
