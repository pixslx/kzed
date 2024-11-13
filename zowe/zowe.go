package zowe

type ZOWE interface {
	JobSubmitDSPath(path string) (ZOWEJobSubmitOutput, error)
	JobSubmitUSSPath(path string) (ZOWEJobSubmitOutput, error)
	JobSubmitJCL(jcl string) (ZOWEJobSubmitOutput, error)
	JobQuery(jobid string) (ZOWEJobSubmitOutput, error)
	JobGetSpoolFiles(jobid string) (ZOWEJobSpoolsOutput, error)

	FilesCreatePDS(input ZOWEDataSetCreateInput) (ZOWEFileOutput, error)
	FilesUploadSTDIN2DS(path string, content string) (ZOWEFileOutput, error)
	FilesDSExists(path string) bool
	FilesDSDelete(path string) error
	FilesDSListMembers(path string) (ZOWEFileOutput, error)
}
