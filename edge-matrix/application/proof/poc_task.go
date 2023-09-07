package proof

import "time"

type PocCpuRequest struct {
	NodeId   string
	Seed     string
	BlockNum uint64
	Start    time.Time
}

type PocCpuData struct {
	Validator string
	Seed      string
	ModelName string
}

type PocSubmitTask struct {
	index int

	// info of the task
	pocSubmitData *PocSubmitData

	// priority of the task (the higher the better)
	priority uint64
}

type PocSubmitData struct {
	ValidationTicket int64
	Validator        string
	Power            int64
	TargetNodeID     string
}

type PocTask struct {
	index int

	// info of the task
	pocCpuData *PocCpuData

	// priority of the task (the higher the better)
	priority uint64
}

// GetPocCpuDataInfo returns the poc information
func (dt *PocTask) GetPocCpuDataInfo() *PocCpuData {
	return dt.pocCpuData
}

// GetPocCpuDataInfo returns the poc information
func (st *PocSubmitTask) GetPocSubmitData() *PocSubmitData {
	return st.pocSubmitData
}
