package cpu

import (
	"testing"
	"time"
)

func TestCPU(t *testing.T) {
	time.Sleep(time.Millisecond * 200)

	for i := 0; i < 6; i++ {
		sys := GetSystemCPU()
		t.Log(sys)

		app := GetProcess()
		t.Log(app)
		time.Sleep(time.Millisecond * 200)
	}
}
