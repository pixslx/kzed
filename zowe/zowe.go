package zowe

type ZOWE interface {
	JobSubmitDSPath(path string) (ZOWEJobSubmitOutput, error)
	JobSubmitUSSPath(path string) (ZOWEJobSubmitOutput, error)
	JobSubmitJCL(jcl string) (ZOWEJobSubmitOutput, error)
	JobQuery(jobid string) (ZOWEJobSubmitOutput, error)
	JobGetSpoolFiles(jobid string) (ZOWEJobSpoolsOutput, error)

	FileUploadSTDIN2DS(content string) (ZOWEFileUploadOutput, error)
}
