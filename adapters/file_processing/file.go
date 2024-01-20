package fileprocessor_adapter

import "mime/multipart"

// FileProcessor is an interface used to get a filename and it's size in bytes
type FileProcessor interface {
	// GetFileNameAndSize will return fileName & fileSize
	// respectively.
	// File size is in bytes
	GetFileNameAndSize(multipart.FileHeader) (string, string)
}
