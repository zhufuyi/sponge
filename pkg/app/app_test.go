package app

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	inits = []Init{
		func() {
			fmt.Println("init config")
		},
	}

	s       = &httpServer{}
	servers = []IServer{s}

	closes = []Close{
		func() error {
			return s.Stop()
		},
	}
)

type httpServer struct{}

func (h *httpServer) Start() error {
	fmt.Println("running http server")
	return nil
}

func (h *httpServer) Stop() error {
	fmt.Println("stop http server")
	return nil
}

func (h *httpServer) String() string {
	return ":8080"
}

func TestNew(t *testing.T) {
	New(inits, servers, closes)
}

func TestApp_Run(t *testing.T) {
	a := New(inits, servers, closes)
	go a.Run()
	time.Sleep(time.Millisecond * 100)
}

func TestApp_stop(t *testing.T) {
	a := New(inits, servers, closes)
	t.Log(a.stop())
}

func TestApp_watch(t *testing.T) {
	a := New(inits, servers, closes)
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	assert.Error(t, a.watch(ctx))
}
