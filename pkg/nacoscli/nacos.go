package nacoscli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// Params nacos参数
type Params struct {
	IPAddr      string // nacos 服务地址
	Port        uint64 // nacos 服务端口
	Scheme      string // http或https
	ContextPath string // path
	NamespaceID string // 名称空间id
	// 如果参数不为空，替换上面和ClientConfig和ServerConfig相同的字段
	clientConfig  *constant.ClientConfig
	serverConfigs []constant.ServerConfig

	Group  string // 分组，dev, prod, test
	DataID string // 配置文件id
	Format string // 配置文件类型: json,yaml,toml
}

func (p *Params) valid() error {
	if p.Group == "" {
		return errors.New("field 'Group' cannot be empty")
	}
	if p.DataID == "" {
		return errors.New("field 'DataID' cannot be empty")
	}
	if p.Format == "" {
		return errors.New("field 'DataID' cannot be empty")
	}
	format := strings.ToLower(p.Format)
	switch format {
	case "json", "yaml", "toml":
		p.Format = format
	case "yml":
		p.Format = "yaml"
	default:
		return fmt.Errorf("config file types 'Format=%s' not supported", p.Format)
	}

	return nil
}

func setParams(params *Params, opts ...Option) {
	o := defaultOptions()
	o.apply(opts...)
	params.clientConfig = o.clientConfig
	params.serverConfigs = o.serverConfigs

	// 创建clientConfig
	if params.clientConfig == nil {
		params.clientConfig = &constant.ClientConfig{
			NamespaceId:         params.NamespaceID,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              os.TempDir() + "/nacos/log",
			CacheDir:            os.TempDir() + "/nacos/cache",
		}
	}

	// 创建serverConfig
	if params.serverConfigs == nil {
		params.serverConfigs = []constant.ServerConfig{
			{
				IpAddr:      params.IPAddr,
				Port:        params.Port,
				Scheme:      params.Scheme,
				ContextPath: params.ContextPath,
			},
		}
	}
}

// Init 从nacos获取配置并解析到struct
func Init(obj interface{}, params *Params, opts ...Option) error {
	err := params.valid()
	if err != nil {
		return err
	}

	setParams(params, opts...)

	// 创建动态配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  params.clientConfig,
			ServerConfigs: params.serverConfigs,
		},
	)
	if err != nil {
		return err
	}

	// 读取配置内容
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: params.DataID,
		Group:  params.Group,
	})
	if err != nil {
		return err
	}

	// 解析配置
	viper.SetConfigType(params.Format)
	err = viper.ReadConfig(bytes.NewBuffer([]byte(content)))
	if err != nil {
		return err
	}
	err = viper.Unmarshal(obj)
	if err != nil {
		return err
	}

	return nil
}

// NewNamingClient 实例化服务注册和发现nacos客户端
func NewNamingClient(nacosIPAddr string, nacosPort int, nacosNamespaceID string, opts ...Option) (naming_client.INamingClient, error) {
	params := &Params{
		IPAddr:      nacosIPAddr,
		Port:        uint64(nacosPort),
		NamespaceID: nacosNamespaceID,
	}
	setParams(params, opts...)

	return clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  params.clientConfig,
			ServerConfigs: params.serverConfigs,
		},
	)
}

//func NewNamingClient(params *Params, opts ...Option) (naming_client.INamingClient, error) {
//	setParams(params, opts...)
//
//	return clients.NewNamingClient(
//		vo.NacosClientParam{
//			ClientConfig:  params.clientConfig,
//			ServerConfigs: params.serverConfigs,
//		},
//	)
//}
