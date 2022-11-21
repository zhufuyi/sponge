package cpu

import "testing"

func TestCgroup(t *testing.T) {
	cgroup := &cgroup{cgroupSet: map[string]string{}}

	us, err := cgroup.CPUCFSQuotaUs()
	t.Log(us, err)

	pus, err := cgroup.CPUCFSPeriodUs()
	t.Log(pus, err)

	usage, err := cgroup.CPUAcctUsage()
	t.Log(usage, err)

	pUsage, err := cgroup.CPUAcctUsagePerCPU()
	t.Log(pUsage, err)

	cpus, err := cgroup.CPUSetCPUs()
	t.Log(cpus, err)

	cgr, err := currentcGroup()
	t.Log(cgr, err)
}

func TestCgroupCPUObj(t *testing.T) {
	cgCPU := &cgroupCPU{
		frequency: 3000,
		quota:     12,
		cores:     12,
		preSystem: 10,
		preTotal:  10,
	}

	usage, err := cgCPU.Usage()
	t.Log(usage, err)

	info := cgCPU.Info()
	t.Log(info)
}

func TestCgroupCPU(t *testing.T) {
	usage, err := systemCPUUsage()
	t.Log(usage, err)

	usage, err = totalCPUUsage()
	t.Log(usage, err)

	usages, err := perCPUUsage()
	t.Log(usages, err)

	sets, err := cpuSets()
	t.Log(sets, err)

	quota, err := cpuQuota()
	t.Log(quota, err)

	period, err := cpuPeriod()
	t.Log(period, err)

	freq := cpuFreq()
	t.Log(freq, err)

	maxFreq := cpuMaxFreq()
	t.Log(maxFreq, err)
}
