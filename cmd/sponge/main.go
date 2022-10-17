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

/*
//go:embed templates/sponge
var microServiceFS embed.FS
*/

func main() {
	rand.Seed(time.Now().UnixNano())

	// 初始化模板，执行命令需要依赖真实文件
	err := generate.Init(generate.TplNameSponge, os.TempDir()+gofile.GetPathDelimiter()+"sponge")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 初始FS化模板，执行命令不需要依赖文件
	//replacer.InitFS(gen.TplNameSponge, "templates/sponge", microServiceFS)

	rootCMD := commands.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
