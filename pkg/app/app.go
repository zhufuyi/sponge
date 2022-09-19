package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhufuyi/sponge/pkg/logger"

	"golang.org/x/sync/errgroup"
)

// Init app initialization
type Init func()

// Close app close
type Close func() error

// IServer http or grpc server interface
type IServer interface {
	Start() error
	Stop() error
	String() string
}

// App servers
type App struct {
	inits   []Init
	servers []IServer
	closes  []Close
}

// New create an app
func New(inits []Init, servers []IServer, closes []Close) *App {
	for _, init := range inits {
		init()
	}

	return &App{
		inits:   inits,
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
			logger.Infof("start up %s", s.String())
			if err := s.Start(); err != nil {
				return err
			}
			return nil
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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case <-ctx.Done(): // service error signals
			_ = a.stop()
			return ctx.Err()
		case s := <-quit: // system notification signal
			logger.Infof("receive a quit signal: %s", s.String())
			if err := a.stop(); err != nil {
				return err
			}
			logger.Infof("stop app successfully")
			return nil
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
