//go:build linux || darwin
// +build linux darwin

package stat

import "syscall"

func init() {
	go func() {
		for {
			select {
			case <-notifyCh:
				syscall.Kill(syscall.Getpid(), syscall.SIGTRAP)
			}
		}
	}()
}
