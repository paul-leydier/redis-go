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
	case Error:
		msg = fmt.Sprintf("-%s\r\n", content)
	case BulkString:
		msg = fmt.Sprintf("$%d\r\n%s\r\n", len(content), content)
	default:
		msg = ""
	}
	return []byte(msg)
}

func RespDecode(msg []byte) (RespType, string) {
	if len(msg) == 0 {
		log.Panicf("cannot decode empty msg")
	}
	msg = bytes.Trim(msg, "\x00")
	switch msg[0] {
	case '+':
		decoded := string(msg[1:])
		return SimpleString, strings.TrimRight(decoded, "\r\n")
	case '-':
		decoded := string(msg[1:])
		return Error, strings.TrimRight(decoded, "\r\n")
	case '$':
		return BulkString, parseBulkString(msg)
	default:
		log.Panicf("unknown msg type identifier %b", msg[0])
	}
	return Error, "should not be here"
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
