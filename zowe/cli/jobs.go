package zowecli

import (
	"encoding/json"
	"io"
	"kzed/m/v2/zowe"
	"os/exec"
)

func (z *ZOWECLI) JobSubmitDSPath(path string) (zowe.ZOWEJobSubmitOutput, error) {
	zoweOut, err := exec.Command("sh", "-c", "zowe jobs submit data-set --rfj '"+path+"'").Output()
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	zoweResponse := zowe.ZOWEJobSubmitOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	return zoweResponse, nil
}

func (z *ZOWECLI) JobSubmitUSSPath(path string) (zowe.ZOWEJobSubmitOutput, error) {
	zoweOut, err := exec.Command("sh", "-c", "zowe jobs submit uss-file --rfj '"+path+"'").Output()
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	zoweResponse := zowe.ZOWEJobSubmitOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	return zoweResponse, nil
}

func (z *ZOWECLI) JobSubmitJCL(jcl string) (zowe.ZOWEJobSubmitOutput, error) {
	zoweCmd := exec.Command("sh", "-c", "zowe jobs submit stdin --rfj")
	zoweStdin, err := zoweCmd.StdinPipe()
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	go func() {
		defer zoweStdin.Close()
		io.WriteString(zoweStdin, jcl)
	}()

	zoweOut, err := zoweCmd.Output()
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}
	zoweResponse := zowe.ZOWEJobSubmitOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	return zoweResponse, nil
}

func (z *ZOWECLI) JobQuery(jobid string) (zowe.ZOWEJobSubmitOutput, error) {
	zoweOut, err := exec.Command("sh", "-c", "zowe jobs vw jsbj --rfj "+jobid).Output()
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}
	zoweResponse := zowe.ZOWEJobSubmitOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEJobSubmitOutput{}, err
	}

	return zoweResponse, nil
}

func (z *ZOWECLI) JobGetSpoolFiles(jobid string) (zowe.ZOWEJobSpoolsOutput, error) {
	zoweOut, err := exec.Command("sh", "-c", "zowe jobs vw asc --rfj "+jobid).Output()
	if err != nil {
		return zowe.ZOWEJobSpoolsOutput{}, err
	}
	zoweSpoolsResponse := zowe.ZOWEJobSpoolsOutput{}
	err = json.Unmarshal([]byte(zoweOut), &zoweSpoolsResponse)
	if err != nil {
		return zowe.ZOWEJobSpoolsOutput{}, err
	}

	return zoweSpoolsResponse, nil
}
