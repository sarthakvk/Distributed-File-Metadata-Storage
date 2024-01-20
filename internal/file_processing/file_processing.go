package fileheaderprocessor

import (
	"fmt"
	"mime/multipart"

	ks_adapter "github.com/sarthakvk/gofilemeta/adapters/keystore_adapter"
)

// FileHeaderProcessor: It's an implementation to process file
type FileHeaderProcessor struct {
	store ks_adapter.AbstractKeyStore
}

// Implements file processing logic
func (f *FileHeaderProcessor) ProcessFile(fileHeader *multipart.FileHeader) error {
	fileName := fileHeader.Filename
	filseSize := fmt.Sprintf("%d", fileHeader.Size)

	err := f.store.Set(fileName, filseSize)

	return err
}

func NewFileHeaderProcessor(store ks_adapter.AbstractKeyStore) *FileHeaderProcessor {
	return &FileHeaderProcessor{store: store}
}
