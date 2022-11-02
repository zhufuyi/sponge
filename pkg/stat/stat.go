package stat

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/stat/cpu"
	"github.com/zhufuyi/sponge/pkg/stat/mem"

	"go.uber.org/zap"
)

var (
	printInfoInterval = time.Minute // minimum 1 second
	zapLog, _         = zap.NewProduction()
)

// Option set the options field.
type Option func(*options)

type options struct{}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithPrintInterval set print interval
func WithPrintInterval(d time.Duration) Option {
	return func(o *options) {
		if d < time.Second {
			return
		}
		printInfoInterval = d
	}
}

// WithLog set zapLog
func WithLog(l *zap.Logger) Option {
	return func(o *options) {
		if l == nil {
			return
		}
		zapLog = l
	}
}

// Init 初始化
func Init(opts ...Option) {
	o := &options{}
	o.apply(opts...)

	go func() {
		printTick := time.NewTicker(printInfoInterval)
		defer printTick.Stop()

		for {
			select {
			case <-printTick.C:
				printUsageInfo()
			}
		}
	}()
}

func printUsageInfo() {
	defer func() { recover() }()

	mSys := mem.GetSystemMemory()
	mProc := mem.GetProcessMemory()
	cSys := cpu.GetSystemCPU()
	cProc := cpu.GetProcess()

	var cors int32
	for _, ci := range cSys.CPUInfo {
		cors += ci.Cores
	}

	sys := system{
		CPUUsage: cSys.UsagePercent,
		CPUCores: cors,
		MemFree:  mSys.Free,
		MemUsage: mSys.UsagePercent,
	}
	proc := process{
		CPUUsage:   cProc.UsagePercent,
		RSS:        cProc.RSS,
		VMS:        cProc.VMS,
		Alloc:      mProc.Alloc,
		TotalAlloc: mProc.TotalAlloc,
		Sys:        mProc.Sys,
		NumGc:      mProc.NumGc,
	}

	zapLog.Info("statistics",
		zap.Any("system", sys),
		zap.Any("process", proc),
	)
}

type system struct {
	CPUUsage float64 `json:"cpu_usage"` // 系统cpu使用率
	CPUCores int32   `json:"cpu_cores"` // cpu核数，多个cpu累加
	MemFree  uint64  `json:"mem_free"`  // 可用物理内存容量，单位(M)
	MemUsage float64 `json:"mem_usage"` // 内存使用率，单位(%)
}

type process struct {
	CPUUsage   float64 `json:"cpu_usage"`   // 进程cpu使用率
	RSS        uint64  `json:"rss"`         // 使用物理内存，单位(M)
	VMS        uint64  `json:"vms"`         // 使用虚拟内存，单位(M)
	Alloc      uint64  `json:"alloc"`       // 分配内存容量，单位(M)
	TotalAlloc uint64  `json:"total_alloc"` // 累计分配内存容量，单位(M)
	Sys        uint64  `json:"sys"`         // 从系统申请内存容量，单位(M)
	NumGc      uint32  `json:"num_gc"`      // 已完成的GC周期的数量
}
