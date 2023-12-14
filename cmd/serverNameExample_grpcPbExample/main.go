// Package main is the grpc server of the application.
package main

import (
	"github.com/zhufuyi/sponge/cmd/serverNameExample_grpcPbExample/initial"

	"github.com/zhufuyi/sponge/pkg/app"
)

func main() {
	initial.InitApp()
	services := initial.CreateServices()
	closes := initial.Close(services)

	a := app.New(services, closes)
	a.Run()
}
