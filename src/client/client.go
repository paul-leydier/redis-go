package redis

import (
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

func (r *Client) Ping(content string) (string, error) {
	_, err := r.conn.Write(core.RespEncode(core.SimpleString, "PING "+content))
	if err != nil {
		return "", err
	}
	encodedResponse := make([]byte, 64)
	_, err = r.conn.Read(encodedResponse)
	resp := core.RespDecode(encodedResponse)
	if err != nil {
		return "", err
	}
	return resp, nil
}
