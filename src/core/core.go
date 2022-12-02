package core

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

type RespType int

const (
	SimpleString RespType = iota
	Error
	Integer
	BulkString
	Array
)

func RespEncode(msgType RespType, content string) []byte {
	var msg string
	switch msgType {
	case SimpleString:
		msg = fmt.Sprintf("+%s\r\n", content)
	default:
		msg = ""
	}
	return []byte(msg)
}

func RespDecode(msg []byte) string {
	if len(msg) == 0 {
		return ""
	}
	msg = bytes.Trim(msg, "\x00")
	switch msg[0] {
	case '+':
		decoded := string(msg[1:])
		return strings.TrimRight(decoded, "\r\n")
	default:
		log.Panicf("unknown msg type identifier %b", msg[0])
	}
	return ""
}
