package fileprocessor_adapter

import "mime/multipart"

// FileProcessor is an interface, it defines a way to process files
type FileProcessor interface {
	// GetFileNameAndSize will return fileName & fileSize
	// respectively.
	// File size is in bytes
	ProcessFile(multipart.FileHeader) error
}
