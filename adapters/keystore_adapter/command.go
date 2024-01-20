package keystore_adapter

import (
	"fmt"
)

type Operation string

const (
	GET           Operation = "GET"
	SET           Operation = "SET"
	DELETE        Operation = "DELETE"
	GET_OR_CREATE Operation = "GET_OR_CREATE"
)

// Commands define the commands used to interact with our
// distributed store
type Command struct {
	Operation Operation `json:"command"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
}

// Error for unknown command
type ErrUnknownCommand struct {
	Operation Operation
}

func (err *ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command: %s", err.Operation)
}
