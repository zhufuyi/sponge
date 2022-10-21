package conf

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var c = make(map[string]interface{})

func init() {
	// result error test
	_ = Parse("test.yml", nil)
	// not found error test
	_ = Parse("notfound.yml", &c)

	err := Parse("test.yml", &c)
	if err != nil {
		panic(err)
	}
}

func TestShow(t *testing.T) {
	t.Log(Show(c))
	t.Log(Show(make(chan string)))
}

func Test_replaceDSN(t *testing.T) {
	dsn := "default:123456@192.168.3.37:6379/0"
	t.Log(replaceDSN(dsn))

	dsn = "default:123456:192.168.3.37:6379/0"
	t.Log(replaceDSN(dsn))
}

func Test_replacePWD(t *testing.T) {
	var keywords []string
	keywords = append(keywords, `"dsn"`, `"password"`, `"name"`)
	str := Show(c, keywords...)

	fmt.Printf(replacePWD(str))
}

// 测试监听配置文件更新
func TestWatch(t *testing.T) {
	time.Sleep(time.Second)
	conf := make(map[string]interface{})

	fs := []func(){
		func() {
			fmt.Println("更新字段1")
		},
		func() {
			fmt.Println("更新字段2")
		},
	}

	err := Parse("test.yml", &conf, fs...)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)
	content, _ := os.ReadFile("test.yml")
	contentChange := append(content, byte('#'))
	time.Sleep(time.Millisecond * 100)
	_ = os.WriteFile("test.yml", contentChange, 0666) // 修改文件
	time.Sleep(time.Millisecond * 100)
	_ = os.WriteFile("test.yml", content, 0666) // 还原文件
	time.Sleep(time.Millisecond * 100)
}
