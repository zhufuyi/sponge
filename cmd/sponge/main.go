package main

import (
	"fmt"
	"os"

	"github.com/zhufuyi/sponge/cmd/sponge/commands"
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

func main() {
	err := generate.Init(generate.TplNameSponge, os.TempDir()+gofile.GetPathDelimiter()+"sponge")
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
