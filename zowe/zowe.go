package zowe

type ZOWE interface {
	JobSubmitPath(path string) (ZOWEJobSubmitOutput, error)
	JobSubmitJCL(jcl string) (ZOWEJobSubmitOutput, error)
	JobQuery(jobid string) (ZOWEJobSubmitOutput, error)
	JobGetSpoolFiles(jobid string) (ZOWEJobSpoolsOutput, error)
}
