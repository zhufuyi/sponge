package nacoscli

import (
	"os"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/stretchr/testify/assert"
)

var (
	ipAddr      = "192.168.3.37"
	port        = 8848
	namespaceID = "3454d2b5-2455-4d0e-bf6d-e033b086bb4c"
)

func TestParse(t *testing.T) {
	// 方式一：
	conf := new(map[string]interface{})
	params := &Params{
		IPAddr:      ipAddr,
		Port:        uint64(port),
		NamespaceID: namespaceID,
		Group:       "dev",
		DataID:      "serverNameExample.yml",
		Format:      "yaml",
	}
	err := Init(conf, params)
	//assert.NoError(t, err)
	t.Log(err, conf)

	// 方式二：直接设置ClientConfig和ServerConfig
	conf = new(map[string]interface{})
	params = &Params{
		Group:  "dev",
		DataID: "serverNameExample.yml",
		Format: "yaml",
	}
	clientConfig := &constant.ClientConfig{
		NamespaceId:         namespaceID,
		TimeoutMs:           1000,
		NotLoadCacheAtStart: true,
		LogDir:              os.TempDir() + "/nacos/log",
		CacheDir:            os.TempDir() + "/nacos/cache",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: ipAddr,
			Port:   uint64(port),
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
	namingClient, err := NewNamingClient(ipAddr, port, namespaceID)
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

	err = Init(nil, p)
	assert.Error(t, err)
}
