package core

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
)

func (r RespElem) String() (string, error) {
	if r.Type == Error {
		return "", errors.New(r.Content.(string))
	}
	if r.Type != SimpleString && r.Type != BulkString {
		return "", fmt.Errorf("invalid RespType: expected %d or %d, got %d", SimpleString, BulkString, r.Type)
	}
	return r.Content.(string), nil
}

func (r RespElem) Int() (int, error) {
	if r.Type == Error {
		return 0, errors.New(r.Content.(string))
	}
	if r.Type != Integer {
		return 0, fmt.Errorf("invalid RespType: expected %d, got %d", Integer, r.Type)
	}
	return r.Content.(int), nil
}

func RespDecode(msg []byte) RespElem {
	if len(msg) == 0 {
		log.Panicf("cannot decode empty msg")
	}
	msg = bytes.Trim(msg, "\x00")
	switch msg[0] {
	case '+':
		decoded := string(msg[1:])
		return RespElem{SimpleString, strings.TrimRight(decoded, "\r\n")}
	case '-':
		decoded := string(msg[1:])
		return RespElem{Error, strings.TrimRight(decoded, "\r\n")}
	case '$':
		return RespElem{BulkString, parseBulkString(msg)}
	default:
		panic(fmt.Sprintf("unknown msg type identifier %b", msg[0]))
	}
}

func parseBulkStringNaive(msg []byte) string {
	parts := bytes.SplitN(msg, []byte("\r\n"), 2)
	if len(parts) < 2 {
		log.Panicf("cannot decode BulkString %b", msg)
	}
	return string(bytes.TrimRight(parts[1], "\r\n"))
}

func parseBulkString(msg []byte) string {
	msgLength := uint8(0)
	prefixLength := 0
	for i := 1; msg[i] != '\r'; i++ {
		msgLength = (msgLength * 10) + (msg[i] - '0')
		prefixLength = i
	}
	return string(msg[prefixLength+3 : prefixLength+3+int(msgLength)])
}
