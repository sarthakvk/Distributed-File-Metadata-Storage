package fileheaderprocessor

import (
	"fmt"
	"mime/multipart"
)

type FileHeaderProcessor int

func (f FileHeaderProcessor) GetFileNameAndSize(fileHeader multipart.FileHeader) (string, string) {
	fileName := fileHeader.Filename
	filseSize := fmt.Sprintf("%d", fileHeader.Size)

	return fileName, filseSize
}

func NewFileHeaderProcessor() FileHeaderProcessor {
	return (FileHeaderProcessor)(1)
}
