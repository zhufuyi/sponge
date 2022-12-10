package prof

import (
	"testing"
	"time"
)

func TestProfile(t *testing.T) {
	EnableTrace()
	SetDurationSecond(2)

	p := NewProfile()

	p.StartOrStop()
	time.Sleep(time.Second)
	p.StartOrStop()
	time.Sleep(time.Second)

	p.StartOrStop()
	time.Sleep(time.Millisecond * 2100)
}
