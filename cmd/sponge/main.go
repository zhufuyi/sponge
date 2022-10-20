package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/zhufuyi/sponge/cmd/sponge/commands"
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := generate.Init(generate.TplNameSponge, os.TempDir()+gofile.GetPathDelimiter()+"sponge")
	if err != nil {
		fmt.Println(err)
		return
	}

	rootCMD := commands.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
