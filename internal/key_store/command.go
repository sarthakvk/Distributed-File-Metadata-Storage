package keystore

import (
	"encoding/json"
	"errors"

	ks_adapter "github.com/sarthakvk/gofilemeta/adapters/keystore_adapter"
)

// GetCommand: creates an Command object from raw byte data
func GetCommand(data []byte) (*ks_adapter.Command, error) {
	var cmd ks_adapter.Command
	err := json.Unmarshal(data, &cmd)

	if err != nil {
		logger.Error("error decoding the command")
		return nil, err
	}

	switch cmd.Operation {
	case ks_adapter.GET, ks_adapter.SET, ks_adapter.DELETE, ks_adapter.GET_OR_CREATE:
		return &cmd, nil
	default:
		logger.Debug("unknown command recieved")
		return nil, &ks_adapter.ErrUnknownCommand{Operation: cmd.Operation}
	}
}

// getRawCommand reverses the function of GetCommand
func getRawCommand(cmd ks_adapter.Command) ([]byte, error) {
	raw_cmd, err := json.Marshal(cmd)

	if err != nil {
		logger.Error("unable to marshal the delete request command!")
		return nil, errors.New("unable to marshal the delete request command")
	}

	return raw_cmd, nil
}
