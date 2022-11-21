package utils

import (
	"context"
	"testing"
	"time"
)

func TestSafeRun(t *testing.T) {
	SafeRun(context.Background(), func(ctx context.Context) {
		t.Log("safe run")
	})

	SafeRun(context.Background(), func(ctx context.Context) {
		panic("run panic")
	})
}

func TestSafeRunWithTimeout(t *testing.T) {
	SafeRunWithTimeout(time.Second, func(cancel context.CancelFunc) {
		t.Log("safe run with timeout")
		cancel()
	})

	SafeRunWithTimeout(time.Second, func(cancel context.CancelFunc) {
		panic("run panic")
	})

	SafeRunWithTimeout(time.Millisecond*100, func(cancel context.CancelFunc) {
		time.Sleep(time.Millisecond * 120)
		cancel()
	})
}
