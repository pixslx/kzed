package zowesdk

import "kzed/m/v2/zowe"

func (z *ZOWESDK) FilesUploadSTDIN2DS(path string, content string) (zowe.ZOWEFileOutput, error) {
	return zowe.ZOWEFileOutput{}, nil
}

func (z *ZOWESDK) FilesCreatePDS(input zowe.ZOWEDataSetCreateInput) (zowe.ZOWEFileOutput, error) {
	return zowe.ZOWEFileOutput{}, nil
}

func (z *ZOWESDK) FilesDSExists(path string) bool {
	return false
}

func (z *ZOWESDK) FilesDSDelete(path string) error {
	return nil
}

func (z *ZOWESDK) FilesDSListMembers(path string) (zowe.ZOWEFileOutput, error) {
	return zowe.ZOWEFileOutput{}, nil
}
