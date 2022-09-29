package nacoscli

import (
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/stretchr/testify/assert"
	"testing"
)

type config struct {
	Env     string `yaml:"env" json:"env"`
	Host    string `yaml:"hostIP" json:"hostIP"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

func TestParse(t *testing.T) {
	// 方式一：
	conf := &config{}
	params := &Params{
		IPAddr:      "192.168.3.37",
		Port:        8848,
		NamespaceID: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataID:      "user-srv.yml",
		Format:      "yaml",
	}
	err := Init(conf, params)
	assert.NoError(t, err)
	t.Log(conf)

	// 方式二：直接设置ClientConfig和ServerConfig
	conf = &config{}
	params = &Params{
		Group:  "dev",
		DataID: "user-srv.yml",
		Format: "yaml",
	}
	clientConfig := &constant.ClientConfig{
		NamespaceId:         "3c715c7a-9e49-4359-8fe6-ff2c67a3a871",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.3.37",
			Port:   8848,
		},
	}
	err = Init(conf, params,
		WithClientConfig(clientConfig),
		WithServerConfigs(serverConfigs),
	)
	assert.NoError(t, err)
	t.Log(conf)
}

func TestNewNamingClient(t *testing.T) {
	params := &Params{
		IPAddr:      "192.168.3.37",
		Port:        8848,
		NamespaceID: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataID:      "user-srv.yml",
		Format:      "yaml",
	}

	namingClient, err := NewNamingClient(params)
	assert.NoError(t, err)
	assert.NotNil(t, namingClient)
}
