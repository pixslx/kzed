package zowecli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kzed/m/v2/zowe"
	"os/exec"
	"strings"
)

func (z *ZOWECLI) FilesUploadSTDIN2DS(path string, content string) (zowe.ZOWEFileOutput, error) {
	zoweCmd := exec.Command("sh", "-c", "zowe files ul stds --rfj '"+path+"'")
	zoweStdin, err := zoweCmd.StdinPipe()
	if err != nil {
		return zowe.ZOWEFileOutput{}, err
	}

	go func() {
		defer zoweStdin.Close()
		io.WriteString(zoweStdin, content)
	}()

	zoweOut, err := zoweCmd.Output()
	if err != nil {
		return zowe.ZOWEFileOutput{}, err
	}
	zoweResponse := zowe.ZOWEFileOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEFileOutput{}, err
	}

	return zoweResponse, nil
}

func (z *ZOWECLI) FilesCreatePDS(input zowe.ZOWEDataSetCreateInput) (zowe.ZOWEFileOutput, error) {
	var zoweCmd strings.Builder
	if input.Like != "" {
		zoweCmd.WriteString("zowe files cre ds --rfj '" + input.Name + "'")
		zoweCmd.WriteString(" --lk '" + input.Like + "'")
	} else {
		zoweCmd.WriteString("zowe files cre pds --rfj '" + input.Name + "'")
	}
	if input.AllocationSpaceUnit != "" {
		zoweCmd.WriteString(" --asu " + input.AllocationSpaceUnit)
	}
	if input.DataClass != "" {
		zoweCmd.WriteString(" --dc " + input.DataClass)
	}
	if input.DataSetType != "" {
		zoweCmd.WriteString(" --dt " + input.DataSetType)
	}
	if input.DeviceType != "" {
		zoweCmd.WriteString(" --dt " + input.DeviceType)
	}
	if input.ManagementClass != "" {
		zoweCmd.WriteString(" --mc " + input.ManagementClass)
	}
	if input.RecordFormat != "" {
		zoweCmd.WriteString(" --rf " + input.RecordFormat)
	}
	if input.Size != "" {
		zoweCmd.WriteString(" --sz " + input.Size)
	}
	if input.StorageClass != "" {
		zoweCmd.WriteString(" --sc " + input.StorageClass)
	}
	if input.VolumeSerial != "" {
		zoweCmd.WriteString(" --vs " + input.VolumeSerial)
	}
	if input.BlockSize != 0 {
		zoweCmd.WriteString(" --bs " + fmt.Sprint(input.BlockSize))
	}
	if input.DirectoryBlocks != 0 {
		zoweCmd.WriteString(" --db " + fmt.Sprint(input.DirectoryBlocks))
	}
	if input.PrimarySpace != 0 {
		zoweCmd.WriteString(" --ps " + fmt.Sprint(input.PrimarySpace))
	}
	if input.RecordLength != 0 {
		zoweCmd.WriteString(" --rl " + fmt.Sprint(input.RecordLength))
	}
	if input.SecondarySpace != 0 {
		zoweCmd.WriteString(" --ss " + fmt.Sprint(input.SecondarySpace))
	}

	z.logger.Info("Issuing ZOWE Command " + zoweCmd.String())
	zoweOut, err := exec.Command("sh", "-c", zoweCmd.String()).Output()
	if err != nil {
		z.logger.Error(err, string(zoweOut))
		return zowe.ZOWEFileOutput{}, err
	}

	zoweResponse := zowe.ZOWEFileOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEFileOutput{}, err
	}

	return zoweResponse, nil
}

func (z *ZOWECLI) FilesDSExists(path string) bool {
	zoweOut, err := exec.Command("sh", "-c", "zowe files ls ds --rfj '"+path+"'").Output()
	if err != nil {
		z.logger.Error(err, string(zoweOut))
		return false
	}

	zoweResponse := zowe.ZOWEFileOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return false
	}

	if zoweResponse.Data.APIResponse.ReturnedRows > 0 {
		return true
	} else {
		return false
	}
}

func (z *ZOWECLI) FilesDSDelete(path string) error {
	zoweOut, err := exec.Command("sh", "-c", "zowe files del ds --rfj '"+path+"' -f").Output()
	if err != nil {
		z.logger.Error(err, string(zoweOut))
		return err
	}

	zoweResponse := zowe.ZOWEFileOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return err
	}

	if zoweResponse.ExitCode > 0 {
		return errors.New("PDS can not be deleted " + path)
	} else {
		return nil
	}
}

func (z *ZOWECLI) FilesDSListMembers(path string) (zowe.ZOWEFileOutput, error) {
	zoweOut, err := exec.Command("sh", "-c", "zowe files ls am --rfj '"+path+"'").Output()
	if err != nil {
		return zowe.ZOWEFileOutput{}, err
	}

	zoweResponse := zowe.ZOWEFileOutput{}
	err = json.Unmarshal(zoweOut, &zoweResponse)
	if err != nil {
		return zowe.ZOWEFileOutput{}, err
	}

	return zoweResponse, nil
}
