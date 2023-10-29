package goredis

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
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
			args:    args{"127.0.0.1:6379"},
			wantErr: false,
		},
		{
			name:    "has password, no db",
			args:    args{"root:123456@127.0.0.1:6379"},
			wantErr: false,
		},
		{
			name:    "no password, has db",
			args:    args{"127.0.0.1:6379/5"},
			wantErr: false,
		},
		{
			name:    "has password, has db",
			args:    args{"root:123456@127.0.0.1:6379/5"},
			wantErr: false,
		},
		{
			name:    "has redis prefix",
			args:    args{"redis://root:123456@127.0.0.1:6379/7"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Init(tt.args.redisURL,
				WithEnableTrace(),
				WithDialTimeout(time.Second),
				WithReadTimeout(time.Second),
				WithWriteTimeout(time.Second),
				WithTLSConfig(&tls.Config{}),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInit2(t *testing.T) {
	rdb := Init2("127.0.0.1:6379", "123456", 0,
		WithEnableTrace(),
		WithDialTimeout(time.Second),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
		WithTLSConfig(&tls.Config{}),
	)
	assert.NotNil(t, rdb)
}

func TestInitSentinel(t *testing.T) {
	addrs := []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382"}
	rdb := InitSentinel("master", addrs, "default", "123456",
		WithEnableTrace(),
		WithDialTimeout(time.Second),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
		WithTLSConfig(&tls.Config{}),
	)
	assert.NotNil(t, rdb)
}

func TestInitCluster(t *testing.T) {
	addrs := []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382"}
	clusterRdb := InitCluster(addrs, "default", "123456",
		WithEnableTrace(),
		WithDialTimeout(time.Second),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
		WithTLSConfig(&tls.Config{}),
	)
	assert.NotNil(t, clusterRdb)
}
