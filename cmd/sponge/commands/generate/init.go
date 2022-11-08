package generate

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/replacer"
)

const warnSymbol = "⚠ "

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Replacers 模板名称对应的接口
var Replacers = map[string]replacer.Replacer{}

// Template 模板信息
type Template struct {
	Name     string
	FS       embed.FS
	FilePath string
}

// Init 初始化模板
func Init(name string, filepath string) error {
	// 判断模板文件是否存在，不存在，提示先更新
	if !gofile.IsExists(filepath) {
		if isShowCommand() {
			return nil
		}
		return fmt.Errorf("%s not yet initialized, run the command 'sponge init'", warnSymbol)
	}

	var err error
	if _, ok := Replacers[name]; ok {
		panic(fmt.Sprintf("template name '%s' already exists", name))
	}
	Replacers[name], err = replacer.New(filepath)
	if err != nil {
		return err
	}

	return nil
}

// InitFS 初始化FS模板
func InitFS(name string, filepath string, fs embed.FS) {
	var err error
	if _, ok := Replacers[name]; ok {
		panic(fmt.Sprintf("template name '%s' already exists", name))
	}
	Replacers[name], err = replacer.NewFS(filepath, fs)
	if err != nil {
		panic(err)
	}
}

func isShowCommand() bool {
	l := len(os.Args)

	// sponge
	if l == 1 {
		return true
	}

	// sponge update or sponge -h
	if l == 2 {
		if os.Args[1] == "init" || os.Args[1] == "-h" {
			return true
		}
		return false
	}
	if l > 2 {
		return strings.Contains(strings.Join(os.Args[:3], ""), "init")
	}

	return false
}
