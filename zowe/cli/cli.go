package zowecli

import (
	"kzed/m/v2/zowe"

	"github.com/go-logr/logr"
)

type ZOWECLI struct {
	logger logr.Logger
}

func NewZOWECLI(logger logr.Logger) zowe.ZOWE {
	return &ZOWECLI{
		logger: logger,
	}
}
