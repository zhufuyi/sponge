package main

import (
	"github.com/zhufuyi/sponge/cmd/serverNameExample_httpExample/initial"

	"github.com/zhufuyi/sponge/pkg/app"
)

// @title serverNameExample api docs
// @description http server api docs
// @schemes http https
// @version v0.0.0
// @host localhost:8080
func main() {
	initial.Config()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
