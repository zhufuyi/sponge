package conf

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)
}

func TestShow(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)
	t.Log(Show(c))
}

func Test_replaceDSN(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)

	str := Show(c)

	fmt.Printf(replaceDSN(str))
}

func Test_replacePWD(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)

	var keywords []string
	keywords = append(keywords, `"dsn"`, `"password"`)
	str := Show(c)

	fmt.Printf(replacePWD(str, keywords...))
}

func Test_watchConfig(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c, func() {
		t.Log("enable watch config file")
	})
	assert.NoError(t, err)

	watchConfig(c)
}

// 测试监听配置文件更新
func TestWatch(t *testing.T) {
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

	for i := 0; i < 1; i++ { // 设置在等待时间内，修改配置文件env字段，查看是否有变化
		fmt.Println(conf["app"])
		time.Sleep(time.Second)
	}
}
