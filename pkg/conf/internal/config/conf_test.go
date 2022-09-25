package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/conf"
)

func TestParseYAML(t *testing.T) {
	err := Init("conf.yml") // 解析yaml文件
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(conf.Show(config))
	fmt.Println()
}

// 测试更新配置文件
func TestWatch(t *testing.T) {
	fs := []func(){
		func() {
			fmt.Println("更新字段1")
		},
		func() {
			fmt.Println("更新字段2")
		},
	}

	err := Init("conf.yml", fs...)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 1; i++ { // 设置100秒等待时间，修改配置文件env字段
		fmt.Println("port:", Get().App.Env)
		time.Sleep(time.Second)
	}
}
