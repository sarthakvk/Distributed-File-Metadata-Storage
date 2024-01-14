package filehandling

import (
	"github.com/sarthakvk/hex-app/ports"
)

type FileHandler struct {
	filename  string
	filesize  int
	logger    ports.Logger
	consensus ports.Consensus
}

func NewFileHandler(filename string, filesize int, consensus ports.Consensus, logger ports.Logger) *FileHandler {
	return &FileHandler{filename: filename, filesize: filesize, logger: logger, consensus: consensus}
}

func (f *FileHandler) SaveDataToDB() {
	f.logger.Log("Init: saving data to DB")
	if f.consensus.GetConsensus() {
		return
	}
}
