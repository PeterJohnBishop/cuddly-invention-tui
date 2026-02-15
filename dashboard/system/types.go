package system

// error
type errMsg struct {
	Err error
}

func (e errMsg) Error() string { return e.Err.Error() }

// CPU
type cpuUsageMsg struct {
	Percent []float64
}

// VirtualMemory
type virtualMemoryMsg struct {
	Total       uint64
	Free        uint64
	UsedPercent float64
}

// Swap Memeory (disk space used as RAM on MacOS)
type swapMemoryMsg struct {
	Total       uint64
	Free        uint64
	UsedPercent float64
}

// Disk
type diskMsg struct {
	Total       uint64
	Free        uint64
	UsedPercent float64
}
