// Serving errors -----------------------------------------

package server

import (
	"fmt"
	"log"
	"redis-go/core"
)

func handleServingError(err error) []byte {
	if icerr, ok := err.(InvalidCommandError); ok {
		return core.RespElem{
			Type:    core.Error,
			Content: fmt.Sprintf("ERR %s", icerr.Error()),
		}.Encode()
	}
	log.Printf("serving error - %s", err)
	return core.RespElem{
		Type:    core.Error,
		Content: "ERR internal error",
	}.Encode()
}

type InvalidCommandError struct {
	received string
}

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf("unknown command '%s'", e.received)
}
