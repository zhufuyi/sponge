// Package mem is a library that counts system and process memory usage.
package mem

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
)

// System memory information
type System struct {
	Total        uint64  `json:"total"`         // total physical memory capacity, unit(M)
	Free         uint64  `json:"free"`          // free physical memory capacity, unit(M)
	UsagePercent float64 `json:"usage_percent"` // memory usage, unit(%)
}

// Process memory information
type Process struct {
	Alloc      uint64 `json:"alloc"`       // allocated memory capacity, unit(M)
	TotalAlloc uint64 `json:"total_alloc"` // cumulative allocated memory capacity, unit(M)
	Sys        uint64 `json:"sys"`         // requesting memory capacity from the system, unit(M)
	NumGc      uint32 `json:"num_gc"`      // number of GC cycles
}

// GetSystemMemory get system memory
func GetSystemMemory() *System {
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("mem.VirtualMemory error: %v\n", err)
		return &System{}
	}

	return &System{
		Total:        info.Total >> 20,
		Free:         info.Free >> 20,
		UsagePercent: info.UsedPercent,
	}
}

// GetProcessMemory get process memory
func GetProcessMemory() *Process {
	info := &runtime.MemStats{}
	runtime.ReadMemStats(info)

	return &Process{
		Alloc:      info.Alloc >> 20,
		TotalAlloc: info.TotalAlloc >> 20,
		Sys:        info.Sys >> 20,
		NumGc:      info.NumGC,
	}
}
