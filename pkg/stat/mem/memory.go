package mem

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
)

// System 系统内存信息
type System struct {
	Total        uint64  `json:"total"`         // 物理内存总容量，单位(M)
	Free         uint64  `json:"free"`          // 可用物理内存容量，单位(M)
	UsagePercent float64 `json:"usage_percent"` // 内存使用率，单位(%)
}

// Process 进程内存信息
type Process struct {
	Alloc      uint64 `json:"alloc"`       // 分配内存容量，单位(M)
	TotalAlloc uint64 `json:"total_alloc"` // 累计分配内存容量，单位(M)
	Sys        uint64 `json:"sys"`         // 从系统申请内存容量，单位(M)
	NumGc      uint32 `json:"num_gc"`      // 已完成的GC周期的数量
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
