package core

import (
	"fmt"
	"log"
)

type RespType int

const (
	SimpleString RespType = iota
	Error
	Integer
	BulkString
	Array
)

type RespElem struct {
	Type    RespType
	Content any
}

func (r RespElem) Encode() []byte {
	var msg string
	switch r.Type {
	case SimpleString:
		msg = fmt.Sprintf("+%s\r\n", r.Content.(string))
	case Error:
		msg = fmt.Sprintf("-%s\r\n", r.Content.(string))
	case Integer:
		msg = fmt.Sprintf(":%d\r\n", r.Content.(int))
	case BulkString:
		c := r.Content.(string)
		msg = fmt.Sprintf("$%d\r\n%s\r\n", len(c), c)
	case Array:
		a := r.Content.([]RespElem)
		msg = encodeArray(a)
	default:
		log.Panicf("unknown RespType %d", r.Type)
	}
	return []byte(msg)
}

func encodeArray(arr []RespElem) string {
	var encodedElements []byte
	for _, a := range arr {
		encodedElements = append(encodedElements, a.Encode()...)
	}
	return fmt.Sprintf("*%d\r\n%s", len(arr), encodedElements)
}
