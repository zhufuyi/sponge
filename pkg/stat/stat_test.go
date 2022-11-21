package stat

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	l, _ := zap.NewDevelopment()
	Init(
		// test empty
		WithLog(nil),
		WithPrintInterval(0),

		WithLog(l),
		WithPrintInterval(time.Second),
	)

	time.Sleep(time.Second * 2)
}
