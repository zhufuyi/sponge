package stat

import (
	"fmt"
	"time"
)

var (
	cpuThreshold            = 0.8
	memoryThreshold         = 0.8
	triggerInterval float64 = 900 // unit(s)
)

// AlarmOption set the alarm options field.
type AlarmOption func(*alarmOptions)

type alarmOptions struct{}

func (o *alarmOptions) apply(opts ...AlarmOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithCPUThreshold set cpu threshold, range 0 to 1
func WithCPUThreshold(threshold float64) AlarmOption {
	return func(o *alarmOptions) {
		if threshold < 0 || threshold >= 1.0 {
			return
		}
		cpuThreshold = threshold
	}
}

// WithMemoryThreshold set memory threshold, range 0 to 1
func WithMemoryThreshold(threshold float64) AlarmOption {
	return func(o *alarmOptions) {
		if threshold < 0 || threshold >= 1.0 {
			return
		}
		memoryThreshold = threshold
	}
}

type statGroup struct {
	data    [3]*statData
	alarmAt time.Time
}

func newStatGroup() *statGroup {
	return &statGroup{data: [3]*statData{}}
}

func (g *statGroup) check(sd *statData) bool {
	if g.data[0] == nil {
		g.data[0] = sd
		return false
	} else if g.data[1] == nil {
		g.data[1] = g.data[0]
		g.data[0] = sd
		return false
	}
	g.data[2] = g.data[1]
	g.data[1] = g.data[0]
	g.data[0] = sd

	if g.checkCPU(cpuThreshold) || g.checkMemory(memoryThreshold) {
		if g.alarmAt.IsZero() {
			g.alarmAt = time.Now()
			return true
		}
		if time.Since(g.alarmAt).Seconds() >= triggerInterval {
			g.alarmAt = time.Now()
			return true
		}
	}

	return false
}

func (g *statGroup) checkCPU(threshold float64) bool {
	if g.data[0].sys.CPUCores == 0 {
		return false
	}

	// average cpu usage exceeds cpuCores*threshold
	average := (g.data[0].proc.CPUUsage + g.data[1].proc.CPUUsage + g.data[2].proc.CPUUsage) / 3 / float64(g.data[0].sys.CPUCores)
	threshold = threshold * 100
	if average >= threshold {
		fmt.Printf("[cpu] processes cpu usage(%.f%%) exceeds %.f%%\n", average, threshold)
		return true
	}

	return false
}

func (g *statGroup) checkMemory(threshold float64) bool {
	if g.data[0].sys.MemTotal == 0 {
		return false
	}

	// processes occupying more than threshold of system memory
	procAverage := (g.data[0].proc.RSS + g.data[1].proc.RSS + g.data[2].proc.RSS) / 3
	procAverageUsage := float64(procAverage) / float64(g.data[0].sys.MemTotal)
	if procAverageUsage >= threshold {
		fmt.Printf("[memory] processes memory usage(%.f%%) exceeds %.f%%\n", procAverageUsage*100, threshold*100)
		return true
	}

	return false
}
