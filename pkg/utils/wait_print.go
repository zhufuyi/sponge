package utils

import (
	"context"
	"fmt"
	"time"
)

// WaitPrinter is a waiting printer.
type WaitPrinter struct {
	ctx            context.Context
	cancel         context.CancelFunc
	printFrequency time.Duration
}

// NewWaitPrinter create a new WaitPrinter instance.
func NewWaitPrinter(interval time.Duration) *WaitPrinter {
	ctx, cancel := context.WithCancel(context.Background())
	if interval < time.Millisecond*100 || interval > time.Second*5 {
		interval = time.Millisecond * 500
	}
	return &WaitPrinter{
		ctx:            ctx,
		cancel:         cancel,
		printFrequency: interval,
	}
}

// LoopPrint start the waiting loop and print the running tip message.
func (p *WaitPrinter) LoopPrint(runningTip string) {
	if p == nil {
		return
	}
	go func() {
		symbols := []string{runningTip + ".", runningTip + "..", runningTip +
			"...", runningTip + "....", runningTip + ".....", runningTip + "......"}
		index := 0
		fmt.Printf("\r%s", symbols[index])

		ticker := time.NewTicker(p.printFrequency)
		defer ticker.Stop()

		for {
			select {
			case <-p.ctx.Done():
				return
			case <-ticker.C:
				index++
				if index >= len(symbols) {
					index = 0
					p.clearCurrentLine()
				}
				fmt.Printf("\r%s", symbols[index])
			}
		}
	}()
}

// StopPrint stop the waiting loop and print the tip message.
func (p *WaitPrinter) StopPrint(tip string) {
	if p == nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	p.cancel()
	p.clearCurrentLine()
	if tip == "" {
		return
	}
	fmt.Println(tip)
}

func (p *WaitPrinter) clearCurrentLine() {
	fmt.Print("\033[2K\r")
}
