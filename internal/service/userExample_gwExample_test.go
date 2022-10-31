package service

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/rpcclient"
)

func TestNewUserExampleServiceClient(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer func() { recover() }()
		rpcclient.NewServerNameExampleRPCConn()
	}()

	time.Sleep(time.Millisecond * 200)
	cli := NewUserExampleServiceClient()
	ctx := context.Background()

	go func() {
		reply, err := cli.Create(ctx, nil)
		t.Log(reply, err)
	}()
	go func() {
		reply, err := cli.DeleteByID(ctx, nil)
		t.Log(reply, err)
	}()
	go func() {
		reply, err := cli.UpdateByID(ctx, nil)
		t.Log(reply, err)
	}()
	go func() {
		reply, err := cli.GetByID(ctx, nil)
		t.Log(reply, err)
	}()
	go func() {
		reply, err := cli.ListByIDs(ctx, nil)
		t.Log(reply, err)
	}()
	go func() {
		reply, err := cli.List(ctx, nil)
		t.Log(reply, err)
	}()

	time.Sleep(time.Second * 15)
}
