// Package app is starting and stopping services gracefully, using golang.org/x/sync/errgroup to ensure that multiple services are started properly at the same time.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/zhufuyi/sponge/pkg/prof"
)

// IServer server interface
type IServer interface {
	Start() error
	Stop() error
	String() string
}

// Close app close
type Close func() error

// App servers
type App struct {
	servers []IServer
	closes  []Close
}

// New create an app
func New(servers []IServer, closes []Close) *App {
	return &App{
		servers: servers,
		closes:  closes,
	}
}

// Run servers
func (a *App) Run() {
	// ctx will be notified whenever an error occurs in one of the goroutines
	eg, ctx := errgroup.WithContext(context.Background())

	// start all servers
	for _, server := range a.servers {
		s := server
		eg.Go(func() error {
			fmt.Println(s.String())
			return s.Start()
		})
	}

	// watch and stop app
	eg.Go(func() error {
		return a.watch(ctx)
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}

// watch the os signal and the ctx signal from the errgroup, and stop the service if either signal is triggered
func (a *App) watch(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTRAP)
	profile := prof.NewProfile()

	for {
		select {
		case <-ctx.Done(): // service error
			_ = a.stop()
			return ctx.Err()

		case sigType := <-sig: // system notification signal
			fmt.Printf("received system notification signal: %s\n", sigType.String())
			switch sigType {
			case syscall.SIGTRAP:
				profile.StartOrStop() // start or stop sampling profile
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
				if err := a.stop(); err != nil {
					return err
				}
				fmt.Println("stop app successfully")
				return nil
			}
		}
	}
}

// stopping services and releasing resources
func (a *App) stop() error {
	for _, closeFn := range a.closes {
		if err := closeFn(); err != nil {
			return err
		}
	}
	return nil
}
