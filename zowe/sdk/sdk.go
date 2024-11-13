package zowesdk

import (
	"kzed/m/v2/zowe"

	"github.com/go-logr/logr"
)

type ZOWESDK struct {
	logger logr.Logger
}

func NewZOWESDK(logger logr.Logger) zowe.ZOWE {
	return &ZOWESDK{
		logger: logger,
	}
}
