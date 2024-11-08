package zowesdk

import (
	"kzed/m/v2/zowe"
)

type ZOWESDK struct{}

func NewZOWESDK() zowe.ZOWE {
	return &ZOWESDK{}
}
