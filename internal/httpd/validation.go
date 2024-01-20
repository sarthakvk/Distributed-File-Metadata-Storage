package httpd

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	key_store "github.com/sarthakvk/gofilemeta/adapters/keystore_adapter"
)

func ValidateKeyStoreCommand(request *http.Request) (*key_store.Command, error) {
	if request.Method != http.MethodPost {
		err := errors.New("bad request")
		return nil, err
	}

	var cmd key_store.Command
	body := request.Body

	err := json.NewDecoder(body).Decode(&cmd)

	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}

	return &cmd, nil
}

func ValidateReplicationRequest(request *http.Request) (*ReplicationRequest, error) {
	if request.Method != http.MethodPost {
		err := errors.New("bad request")
		return nil, err
	}

	var req ReplicationRequest
	body := request.Body

	err := json.NewDecoder(body).Decode(&req)

	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}

	return &req, nil
}

func ValidateFileUpload(request *http.Request) (*multipart.FileHeader, error) {
	// Check if the request is a multipart form (indicating file upload)
	if request.Method != http.MethodPost || !strings.HasPrefix(request.Header.Get("Content-Type"), "multipart/form-data") {
		err := errors.New("invalid request. Please send a file")
		return nil, err
	}

	// Parse the form to extract file
	err := request.ParseMultipartForm(10 << 20) // 10 MB limit for the entire form
	if err != nil {
		return nil, err
	}

	// Check if the request contains any files
	if request.MultipartForm == nil || len(request.MultipartForm.File) == 0 {
		err = errors.New("no file uploaded. Please send a file")
		return nil, err
	} else if len(request.MultipartForm.File) != 1 {
		err = errors.New("multiple files not supported")
		return nil, err
	}

	var fileHeader *multipart.FileHeader
	// Iterate over the files in the map
	for _, fileHeaders := range request.MultipartForm.File {
		// Check if only one file is uploaded for each field
		if len(fileHeaders) != 1 {
			err = errors.New("multiple files with same name")
			return nil, err
		}

		fileHeader = fileHeaders[0]
	}
	return fileHeader, err
}
