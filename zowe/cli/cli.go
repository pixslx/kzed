package zowecli

import (
	"kzed/m/v2/zowe"
)

type ZOWECLI struct{}

func NewZOWECLI() zowe.ZOWE {
	return &ZOWECLI{}
}
