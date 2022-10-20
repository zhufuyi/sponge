package nacoscli

import (
	"os"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/stretchr/testify/assert"
)

var addr = "127.0.0.1"
var port uint64 = 8848

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
		IPAddr:      addr,
		Port:        port,
		NamespaceID: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataID:      "user-srv.yml",
		Format:      "yaml",
	}
	err := Init(conf, params)
	//assert.NoError(t, err)
	t.Log(err, conf)

	// 方式二：直接设置ClientConfig和ServerConfig
	conf = &config{}
	params = &Params{
		Group:  "dev",
		DataID: "user-srv.yml",
		Format: "yaml",
	}
	clientConfig := &constant.ClientConfig{
		NamespaceId:         "3c715c7a-9e49-4359-8fe6-ff2c67a3a871",
		TimeoutMs:           1000,
		NotLoadCacheAtStart: true,
		LogDir:              os.TempDir() + "/nacos/log",
		CacheDir:            os.TempDir() + "/nacos/cache",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: addr,
			Port:   port,
		},
	}
	err = Init(conf, params,
		WithClientConfig(clientConfig),
		WithServerConfigs(serverConfigs),
	)
	//assert.NoError(t, err)
	t.Log(err, conf)
}

func TestNewNamingClient(t *testing.T) {
	params := &Params{
		IPAddr:      addr,
		Port:        port,
		NamespaceID: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataID:      "user-srv.yml",
		Format:      "yaml",
	}

	namingClient, err := NewNamingClient(params)
	//assert.NoError(t, err)
	//assert.NotNil(t, namingClient)
	t.Log(err, namingClient)
}

func TestError(t *testing.T) {
	p := &Params{}
	p.Group = ""
	err := p.valid()
	assert.Error(t, err)

	p.Group = "group"
	p.DataID = ""
	err = p.valid()
	assert.Error(t, err)

	p.Group = "group"
	p.DataID = "id"
	p.Format = ""
	err = p.valid()
	assert.Error(t, err)

	p.Group = "group"
	p.DataID = "id"
	p.Format = "yml"
	err = p.valid()
	assert.NoError(t, err)

	p.Group = "group"
	p.DataID = "id"
	p.Format = "unknown"
	err = p.valid()
	assert.Error(t, err)

	err = setParams(p)
	assert.Error(t, err)

	err = Init(nil, p)
	assert.Error(t, err)

	_, err = NewNamingClient(p)
	assert.Error(t, err)
}
