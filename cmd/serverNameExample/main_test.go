package main

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/pkg/app"
)

var chn = make(chan func())

func Test_Main(t *testing.T) {
	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		select {
	//		case testFn, ok := <-chn:
	//			if !ok {
	//				return
	//			}
	//			testFn()
	//		}
	//	}
	//}()

	go func() {
		time.Sleep(time.Second)
		chn <- testMain
		time.Sleep(time.Second)
		chn <- testRegisterCloses
		time.Sleep(time.Second)
		chn <- testGrpcOptions
		time.Sleep(time.Second * 2)
		close(chn)
	}()

	for {
		select {
		case testFn, ok := <-chn:
			if !ok {
				return
			}
			testFn()
		case <-time.After(time.Second * 15):
			return
		}
	}
}

func initLocalConfig() {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	config.Get().App.EnableTracing = true
	config.Get().App.EnableRegistryDiscovery = true
}

var testRegisterCloses = func() {
	initLocalConfig()
	closeFns := registerCloses([]app.IServer{&srv{}})
	for _, fn := range closeFns {
		_ = fn()
	}
}

var testGrpcOptions = func() {
	initLocalConfig()
	defer func() { recover() }()
	grpcOptions()
}

var testMain = func() {
	defer func() {
		recover()
	}()

	main()
}

type srv struct{}

func (s srv) Start() error {
	return nil
}

func (s srv) Stop() error {
	return nil
}

func (s srv) String() string {
	return "foo"
}
