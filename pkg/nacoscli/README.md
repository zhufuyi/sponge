## nacoscli

Get the configuration from the nacos configuration center and parse it into a structure.

### Example of use

```go
	import "github.com/zhufuyi/sponge/pkg/nacoscli"

	// Way 1: Setting parameters
	a := &config{}
	params := &Params{
		IpAddr:      "192.168.3.37",
		Port:        8848,
		NamespaceId: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataId:      "user-srv.yml",
		Format:      "yaml",
	}
	err := nacoscli.Init(a, params)

	// Way 2: Setting up ClientConfig and ServerConfig
	a = &config{}
	params = &Params{
		Group:  "dev",
		DataId: "user-srv.yml",
		Format: "yaml",
	}
	clientConfig := &constant.ClientConfig{
		NamespaceId:         "de7b176e-91cd-49a3-ac83-beb725979775",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              os.TempDir() + "/nacos/log",
		CacheDir:            os.TempDir() + "/nacos/cache",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.3.37",
			Port:   8848,
		},
	}
	err = nacoscli.Init(a, params,
		WithClientConfig(clientConfig),
		WithServerConfigs(serverConfigs),
	)
```
