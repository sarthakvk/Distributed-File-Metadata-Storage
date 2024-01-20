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

type Command struct {
	Operation Operation `json:"command"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
}

type ErrUnknownCommand struct {
	Operation Operation
}

func (err *ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command: %s", err.Operation)
}
