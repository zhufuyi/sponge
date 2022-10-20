## nacoscli

从nacos配置中心获取配置并解析到结构体。

### 使用示例

```go
	// 方式一：设置参数
	a := &config{}
	params := &Params{
		IpAddr:      "192.168.3.37",
		Port:        8848,
		NamespaceId: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataId:      "user-srv.yml",
		Format:      "yaml",
	}
	err := Init(a, params)

	// 方式二：设置ClientConfig和ServerConfig
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
	err = Init(a, params,
		WithClientConfig(clientConfig),
		WithServerConfigs(serverConfigs),
	)
```
