package core

import "fmt"

func RespEncode(responseType string, content string) []byte {
	switch responseType {
	case "string":
		return []byte(fmt.Sprintf("+%s\r\n", content))
	default:
		return []byte("")
	}
}
