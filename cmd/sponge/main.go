// Package main sponge is a powerful tool for generating web and microservice code, a microservice
// framework based on gin and grpc encapsulation, and an open source framework for rapid application
// development. Sponge has a wealth of code generation commands, sponge generate code unified in
// the UI interface operation, it is easy to build a complete project engineering code.
package main

import (
	"fmt"
	"os"

	"github.com/zhufuyi/sponge/cmd/sponge/commands"
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

func main() {
	err := generate.Init(generate.TplNameSponge, commands.GetSpongeDir()+gofile.GetPathDelimiter()+".sponge")
	if err != nil {
		fmt.Printf("\n    %v\n\n", err)
		return
	}

	rootCMD := commands.NewRootCMD()
	if err = rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
