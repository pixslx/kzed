package zowesdk

import (
	"kzed/m/v2/zowe"
)

func (z *ZOWESDK) JobSubmitPath(path string) (zowe.ZOWEJobSubmitOutput, error) {
	return zowe.ZOWEJobSubmitOutput{}, nil
}

func (z *ZOWESDK) JobSubmitJCL(jcl string) (zowe.ZOWEJobSubmitOutput, error) {
	return zowe.ZOWEJobSubmitOutput{}, nil
}

func (z *ZOWESDK) JobQuery(jobid string) (zowe.ZOWEJobSubmitOutput, error) {
	return zowe.ZOWEJobSubmitOutput{}, nil
}

func (z *ZOWESDK) JobGetSpoolFiles(jobid string) (zowe.ZOWEJobSpoolsOutput, error) {
	return zowe.ZOWEJobSpoolsOutput{}, nil
}
