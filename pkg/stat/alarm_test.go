package stat

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_statGroup_check(t *testing.T) {
	testData := `[
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":0.1,"rss":58,"vms":1475,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":0.1,"rss":4500,"vms":5000,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":0.1,"rss":5500,"vms":5000,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":0.1,"rss":5100,"vms":5000,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":0.1,"rss":58,"vms":1475,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":161,"rss":52,"vms":5000,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":158,"rss":53,"vms":5000,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":175,"rss":53,"vms":5000,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},

{"system": {"cpu_usage":0,"cpu_cores":2,"mem_total":0,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":17.5,"rss":53,"vms":50,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}},
{"system": {"cpu_usage":0,"cpu_cores":0,"mem_total":6000,"mem_free":555,"mem_usage":56.7}, "process": {"cpu_usage":17.5,"rss":53,"vms":50,"alloc":14,"total_alloc":27,"sys":39,"num_gc":6}}
]`

	type stData struct {
		System  system  `json:"system"`
		Process process `json:"process"`
	}

	sd := []stData{}
	err := json.Unmarshal([]byte(testData), &sd)
	if err != nil {
		t.Error(err)
		return
	}

	sg := newStatGroup()
	triggerInterval = 1
	for _, data := range sd {
		isAlarm := sg.check(&statData{
			sys:  data.System,
			proc: data.Process,
		})
		t.Log(isAlarm)
		time.Sleep(250 * time.Millisecond)
	}

	sg.checkCPU(0)
}

func Test_alarmOptions_apply(t *testing.T) {
	t.Log(cpuThreshold, memoryThreshold)
	ao := &alarmOptions{}
	ao.apply(
		WithCPUThreshold(-0.5),   // invalid value
		WithMemoryThreshold(1.5), // invalid value

		WithCPUThreshold(0.9),
		WithMemoryThreshold(0.85),
	)
	t.Log(cpuThreshold, memoryThreshold)
}
