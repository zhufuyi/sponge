package mem

import (
	"testing"
	"time"
)

func TestGetSystemMemory(t *testing.T) {
	time.Sleep(time.Millisecond * 200)
	for i := 0; i < 10; i++ {
		sm := GetSystemMemory()
		am := GetProcessMemory()
		t.Log(sm, am)
		time.Sleep(time.Millisecond * 300)
	}
}
