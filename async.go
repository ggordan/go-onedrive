package onedrive

// AsyncJob stores the location (URL) which can be pinged with CheckStatus() to
// check progress of an Async job.
type AsyncJob struct {
	*OneDrive
	Location string
}

// AsyncJobStatus provides information on the status of a asynchronous job progress.
type AsyncJobStatus struct {
	Operation          string  `json:"operation"`
	PercentageComplete float64 `json:"percentageComplete"`
	Status             string  `json:"status"`
}

// CheckStatus returns a new AsyncJobStatus
func (aj AsyncJob) CheckStatus() (*AsyncJobStatus, error) {
	req, err := aj.newRequest("GET", aj.Location, nil, nil)
	if err != nil {
		return nil, err
	}

	ajs := new(AsyncJobStatus)
	_, err = aj.do(req, ajs)
	if err != nil {
		return nil, err
	}

	return ajs, nil
}
