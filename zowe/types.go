package zowe

type ZOWEJobSubmitOutput struct {
	Success  bool                    `json:"success"`
	ExitCode int                     `json:"exitCode"`
	Message  string                  `json:"message"`
	Stdout   string                  `json:"stdout"`
	Stderr   string                  `json:"stderr"`
	Data     ZOWEJobSubmitOutputData `json:"data"`
}

type ZOWEJobSubmitOutputData struct {
	Owner         string `json:"owner"`
	Phase         int    `json:"phase"`
	Subsystem     string `json:"subsystem"`
	PhaseName     string `json:"phase-name"`
	JobCorrelator string `json:"job-correlator"`
	Type          string `json:"type"`
	URL           string `json:"url"`
	Jobid         string `json:"jobid"`
	Class         string `json:"class"`
	FilesURL      string `json:"files-url"`
	Jobname       string `json:"jobname"`
	Status        string `json:"status"`
	Retcode       string `json:"retcode"`
}

type ZOWEJobSpoolsOutput struct {
	Success  bool                      `json:"success"`
	ExitCode int                       `json:"exitCode"`
	Message  string                    `json:"message"`
	Stdout   string                    `json:"stdout"`
	Stderr   string                    `json:"stderr"`
	Data     []ZOWEJobSpoolsOutputData `json:"data"`
}

type ZOWEJobSpoolsOutputData struct {
	ID       int    `json:"id"`
	DdName   string `json:"ddName"`
	StepName string `json:"stepName"`
	ProcName string `json:"procName"`
	Data     string `json:"data"`
}

type ZOWEFileOutput struct {
	Success  bool               `json:"success"`
	ExitCode int                `json:"exitCode"`
	Message  string             `json:"message"`
	Stdout   string             `json:"stdout"`
	Stderr   string             `json:"stderr"`
	Data     ZOWEFileOutputData `json:"data"`
	Error    struct{}           `json:"error"`
}

type ZOWEFileOutputData struct {
	Success         bool   `json:"success"`
	CommandResponse string `json:"commandResponse"`
	APIResponse     struct {
		Items []struct {
			Member string `json:"member"`
		} `json:"items"`
		ReturnedRows int `json:"returnedRows"`
		JSONversion  int `json:"JSONversion"`
	} `json:"apiResponse"`
}

type ZOWEDataSetCreateInput struct {
	Name                string
	AllocationSpaceUnit string
	BlockSize           int
	DataClass           string
	DataSetType         string
	DeviceType          string
	DirectoryBlocks     int
	ManagementClass     string
	PrimarySpace        int
	RecordFormat        string
	RecordLength        int
	SecondarySpace      int
	Size                string
	StorageClass        string
	VolumeSerial        string
	Like                string
}
