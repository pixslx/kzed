package zowecli

import (
	"encoding/json"
	"io"
	"kzed/m/v2/zowe"
	"os/exec"
)

func (z *ZOWECLI) FileUploadSTDIN2DS(content string) (zowe.ZOWEFileUploadOutput, error) {
	zoweCmd := exec.Command("sh", "-c", "zowe files ul stds --rfj")
	zoweStdin, err := zoweCmd.StdinPipe()
	if err != nil {
		return zowe.ZOWEFileUploadOutput{}, err
	}

	go func() {
		defer zoweStdin.Close()
		io.WriteString(zoweStdin, content)
	}()

	zoweOut, err := zoweCmd.Output()
	if err != nil {
		return zowe.ZOWEFileUploadOutput{}, err
	}
	zoweResponse := zowe.ZOWEFileUploadOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEFileUploadOutput{}, err
	}

	return zoweResponse, nil
}
