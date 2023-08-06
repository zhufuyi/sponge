// Package cpu is a library that calculates cpu and memory usage.
package cpu

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type psutilCPU struct {
	interval time.Duration
}

func newPsutilCPU(interval time.Duration) (*psutilCPU, error) {
	psCPU := &psutilCPU{interval: interval}
	_, err := psCPU.Usage()
	if err != nil {
		return nil, err
	}
	return psCPU, nil
}

func (ps *psutilCPU) Usage() (uint64, error) {
	var u uint64
	percents, err := cpu.Percent(ps.interval, false)
	if err == nil {
		u = uint64(percents[0] * 10)
	}
	return u, err
}

func (ps *psutilCPU) Info() Info {
	stats, err := cpu.Info()
	if err != nil {
		return Info{}
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		return Info{}
	}

	return Info{
		Frequency: uint64(stats[0].Mhz),
		Quota:     float64(cores),
	}
}
