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

// Init initialize statistical information
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
	CPUUsage float64 `json:"cpu_usage"` // system cpu usage, unit(%)
	CPUCores int32   `json:"cpu_cores"` // cpu cores, multiple cpu accumulation
	MemFree  uint64  `json:"mem_free"`  // system free physical memory, unit(M)
	MemUsage float64 `json:"mem_usage"` // system memory usage, unit(%)
}

type process struct {
	CPUUsage   float64 `json:"cpu_usage"`   // process cpu usage, unit(%)
	RSS        uint64  `json:"rss"`         // use of physical memory, unit(M)
	VMS        uint64  `json:"vms"`         // use of virtual memory, unit(M)
	Alloc      uint64  `json:"alloc"`       // allocated memory capacity, unit(M)
	TotalAlloc uint64  `json:"total_alloc"` // cumulative allocated memory capacity, unit(M)
	Sys        uint64  `json:"sys"`         // requesting memory capacity from the system, unit(M)
	NumGc      uint32  `json:"num_gc"`      // number of GC cycles
}
