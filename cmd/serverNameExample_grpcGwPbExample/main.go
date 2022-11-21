package main

import (
	"github.com/zhufuyi/sponge/cmd/serverNameExample_grpcGwPbExample/initial"

	"github.com/zhufuyi/sponge/pkg/app"
)

func main() {
	initial.Config()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
