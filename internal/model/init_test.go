package model

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

	config.Get().App.EnableTracing = true
	config.Get().Mysql.EnableLog = true
	gdb := GetDB()
	assert.NotNil(t, gdb)
	time.Sleep(time.Millisecond * 10)
	err = CloseMysql()
	assert.NoError(t, err)
}

func TestInitMysqlError(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}

	// change config error test
	config.Get().Mysql.Dsn = "root:123456@(127.0.0.1:3306)/test"
	_ = CloseMysql()
	InitMysql()
}

func TestCloseMysql(t *testing.T) {
	defer func() { recover() }()
	db = &gorm.DB{}
	err := CloseMysql()
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

	// change config error test
	config.Get().App.EnableTracing = true
	config.Get().Redis.Dsn = "default:123456@127.0.0.1:6379/0"
	redisCli = nil
	_ = CloseRedis()
	_ = GetRedisCli()
}

func TestTableName(t *testing.T) {
	t.Log(new(UserExample).TableName())
}
