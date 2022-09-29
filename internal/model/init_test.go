package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/serverNameExample/config"
)

// 测试时需要连接真实数据

func TestInitMysql(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}

	InitMysql()
	gdb := GetDB()
	assert.NotNil(t, gdb)
	time.Sleep(time.Millisecond * 10)
	err = CloseMysql()
	assert.NoError(t, err)
}

func TestInitRedis(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect redis error info")
		}
	}()

	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}

	InitRedis()
	cli := GetRedisCli()
	assert.NotNil(t, cli)
	time.Sleep(time.Millisecond * 10)
	err = CloseRedis()
	assert.NoError(t, err)
}

func TestTableName(t *testing.T) {
	t.Log(new(UserExample).TableName())
}
