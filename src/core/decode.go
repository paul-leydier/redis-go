package core

import (
	"errors"
	"fmt"
	"log"
)

type EncodedRespElem struct {
	msg    []byte
	cursor int
}

func NewEncodedRespElem(msg []byte) EncodedRespElem {
	return EncodedRespElem{
		msg:    msg,
		cursor: 0,
	}
}

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

func (r RespElem) Array() ([]RespElem, error) {
	if r.Type == Error {
		return nil, errors.New(r.Content.(string))
	}
	if r.Type != Array {
		return nil, fmt.Errorf("invalid RespType: expected %d, got %d", Integer, r.Type)
	}
	return r.Content.([]RespElem), nil
}

func RespDecode(encoded *EncodedRespElem) RespElem {
	if len(encoded.msg) == 0 {
		log.Panicf("cannot decode empty msg")
	}
	//encoded = bytes.Trim(encoded, "\x00")
	switch encoded.msg[encoded.cursor] {
	case '+':
		return decodeSimpleString(encoded)
	case '-':
		return decodeError(encoded)
	case '$':
		return parseBulkString(encoded)
	case '*':
		return decodeArray(encoded)
	default:
		panic(fmt.Sprintf("unknown msg type identifier %b", encoded.msg[encoded.cursor]))
	}
}

func decodeSimpleString(encoded *EncodedRespElem) RespElem {
	encoded.cursor++ // consume '+'
	start := encoded.cursor
	for ; encoded.msg[encoded.cursor] != '\r'; encoded.cursor++ {
	}
	encoded.cursor += 2 // consume '\r\n'
	return RespElem{
		Type:    SimpleString,
		Content: string(encoded.msg[start : encoded.cursor-2]),
	}
}

func decodeError(encoded *EncodedRespElem) RespElem {
	encoded.cursor++ // consume '-'
	start := encoded.cursor
	for ; encoded.msg[encoded.cursor] != '\r'; encoded.cursor++ {
	}
	encoded.cursor += 2 // consume '\r\n'
	return RespElem{
		Type:    Error,
		Content: string(encoded.msg[start : encoded.cursor-2]),
	}
}

func parseBulkString(encoded *EncodedRespElem) RespElem {
	encoded.cursor++ // consume '$'
	msgLength := parsePrefix(encoded)
	start := encoded.cursor
	encoded.cursor += msgLength + 2 // consume message and '\r\n'
	return RespElem{
		Type:    BulkString,
		Content: string(encoded.msg[start : encoded.cursor-2]),
	}
}

func parsePrefix(encoded *EncodedRespElem) int {
	msgLength := 0
	for ; encoded.msg[encoded.cursor] != '\r'; encoded.cursor++ {
		msgLength = (msgLength * 10) + (int(encoded.msg[encoded.cursor]) - '0')
	}
	encoded.cursor += 2 // consume '\r\n'
	return msgLength
}

func decodeArray(encoded *EncodedRespElem) RespElem {
	encoded.cursor++ // consume '*'
	msgLength := parsePrefix(encoded)
	arr := make([]RespElem, msgLength)
	for i := 0; i < msgLength; i++ {
		arr[i] = RespDecode(encoded)
	}
	return RespElem{
		Type:    Array,
		Content: arr,
	}
}
