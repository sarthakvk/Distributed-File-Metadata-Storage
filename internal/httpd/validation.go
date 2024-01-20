package httpd

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	fileprocessing "github.com/sarthakvk/hex-app/adapters/file_processing"
	key_store "github.com/sarthakvk/hex-app/adapters/keystore_adapter"
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

func ValidateFileUpload(request *http.Request, fileProcessor fileprocessing.FileProcessor) (string, string, error) {
	// Check if the request is a multipart form (indicating file upload)
	if request.Method != http.MethodPost || !strings.HasPrefix(request.Header.Get("Content-Type"), "multipart/form-data") {
		err := errors.New("invalid request. Please send a file")
		return "", "", err
	}

	// Parse the form to extract file
	err := request.ParseMultipartForm(10 << 20) // 10 MB limit for the entire form
	if err != nil {
		return "", "", err
	}

	// Check if the request contains any files
	if request.MultipartForm == nil || len(request.MultipartForm.File) == 0 {
		err := errors.New("no file uploaded. Please send a file")
		return "", "", err
	} else if len(request.MultipartForm.File) != 1 {
		err := errors.New("multiple files not supported")
		return "", "", err
	}

	var fileSize, fileName string

	// Iterate over the files in the map
	for _, fileHeaders := range request.MultipartForm.File {
		// Check if only one file is uploaded for each field
		if len(fileHeaders) != 1 {
			err := errors.New("multiple files with same name")
			return "", "", err
		}

		// Access the uploaded file name and size for the single file
		fileHeader := fileHeaders[0]
		fileName, fileSize = fileProcessor.GetFileNameAndSize(*fileHeader)
	}

	return fileName, fileSize, nil
}
