package logging

import (
	"sync"

	"github.com/hashicorp/go-hclog"
)

var (
	logger hclog.Logger
	once   sync.Once

	httpd_logger hclog.Logger
	httpd_once   sync.Once
)

func GetLogger() hclog.Logger {
	once.Do(func() {
		opts := &hclog.LoggerOptions{
			Name:  "raft",
			Level: 1,
		}
		logger = hclog.New(opts)

	})
	return logger
}

func GetHttpdLogger() hclog.Logger {
	httpd_once.Do(func() {
		opts := &hclog.LoggerOptions{
			Name:  "httpd",
			Level: 1,
		}
		httpd_logger = hclog.New(opts)

	})
	return httpd_logger
}
