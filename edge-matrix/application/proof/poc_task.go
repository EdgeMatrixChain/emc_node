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
