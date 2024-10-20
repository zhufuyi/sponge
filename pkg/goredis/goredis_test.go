package goredis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	redisServer, _ := miniredis.Run()
	defer redisServer.Close()
	addr := redisServer.Addr()

	type args struct {
		redisURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    " no password, no db",
			args:    args{addr},
			wantErr: false,
		},
		{
			name:    "has password, no db",
			args:    args{"root:123456@" + addr},
			wantErr: false,
		},
		{
			name:    "no password, has db",
			args:    args{addr + "/5"},
			wantErr: false,
		},
		{
			name:    "has password, has db",
			args:    args{fmt.Sprintf("root:123456@%s/5", addr)},
			wantErr: false,
		},
		{
			name:    "has redis prefix",
			args:    args{fmt.Sprintf("redis://root:123456@%s/5", addr)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, err := Init(tt.args.redisURL,
				WithDialTimeout(time.Second),
				WithReadTimeout(time.Second),
				WithWriteTimeout(time.Second),
				WithEnableTrace(),
				WithTracing(nil),   // nil means no set field
				WithTLSConfig(nil), // nil means no set field
			)
			if (err != nil) != tt.wantErr {
				t.Logf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, rdb)
		})
	}
}

func TestInitSingle(t *testing.T) {
	redisServer, _ := miniredis.Run()
	defer redisServer.Close()
	addr := redisServer.Addr()

	rdb, err := InitSingle(addr, "", 0,
		WithDialTimeout(time.Second),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
		WithTracing(nil),       // nil means no set field
		WithTLSConfig(nil),     // nil means no set field
		WithSingleOptions(nil), // nil means no set field
	)
	assert.Nil(t, err)
	assert.NotNil(t, rdb)
}

func TestInitSentinel(t *testing.T) {
	redisServer, _ := miniredis.Run()
	defer redisServer.Close()
	addr := redisServer.Addr()

	rdb, err := InitSentinel("mymaster", []string{addr}, "", "",
		WithDialTimeout(time.Second),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
		WithTracing(nil),         // nil means no set field
		WithTLSConfig(nil),       // nil means no set field
		WithSentinelOptions(nil), // nil means no set field
	)
	t.Log(err)
	assert.NotNil(t, rdb)
}

func TestInitCluster(t *testing.T) {
	redisServer, _ := miniredis.Run()
	defer redisServer.Close()
	addr := redisServer.Addr()

	clusterRdb, err := InitCluster([]string{addr}, "", "",
		WithDialTimeout(time.Second*15),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
		WithTracing(nil),        // nil means no set field
		WithTLSConfig(nil),      // nil means no set field
		WithClusterOptions(nil), // nil means no set field
	)
	assert.Nil(t, err)
	assert.NotNil(t, clusterRdb)
}
