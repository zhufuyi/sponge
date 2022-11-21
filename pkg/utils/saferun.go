package utils

import (
	"context"
	"fmt"
	"time"
)

// SafeRun safe run
func SafeRun(ctx context.Context, fn func(ctx context.Context)) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	fn(ctx)
}

// SafeRunWithTimeout safe run with limit timeouts
func SafeRunWithTimeout(d time.Duration, fn func(cancel context.CancelFunc)) {
	ctx, cancel := context.WithTimeout(context.Background(), d)

	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
			}
		}()

		fn(cancel)
	}()

	for range ctx.Done() {
		return
	}
}
