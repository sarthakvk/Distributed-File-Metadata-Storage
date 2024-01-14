package hclogadapter

import (
	"github.com/hashicorp/go-hclog"
)

type HcLogAdapter struct {
	backend hclog.Logger
}

var Logger = GetLogger()

func GetLogger() *HcLogAdapter {
	return &HcLogAdapter{backend: hclog.Default()}
}

func (log *HcLogAdapter) Log(msg string) {
	log.backend.Info(msg)
}

func (log *HcLogAdapter) Exception(msg string, err error) {
	log.backend.Error(msg)
	log.backend.Error(err.Error())
}
