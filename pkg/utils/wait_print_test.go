package utils

import (
	"testing"
	"time"
)

func TestLoopPrint(t *testing.T) {
	runningTip := "upgrading sponge binary "
	finishTip := "upgrade sponge binary successfully"
	failedTip := "failed to upgrade sponge binary "

	p := NewWaitPrinter(time.Millisecond * 100)
	p.LoopPrint(runningTip)
	time.Sleep(time.Millisecond * 1000)
	p.StopPrint(finishTip)

	p = NewWaitPrinter(0)
	p.LoopPrint(runningTip)
	time.Sleep(time.Millisecond * 1100)
	p.StopPrint(failedTip)
	time.Sleep(time.Millisecond * 100)
}
