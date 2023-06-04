package proof

type PocCpuData struct {
	NodeId   string
	Seed     string
	BlockNum uint64
}

type PocTask struct {
	index int

	// info of the task
	pocCpuData *PocCpuData

	// priority of the task (the higher the better)
	priority uint64
}

// GetAddrInfo returns the peer information associated with the dial
func (dt *PocTask) GetPocCpuDataInfo() *PocCpuData {
	return dt.pocCpuData
}
