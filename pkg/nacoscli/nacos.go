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

// Params nacos parameters
type Params struct {
	IPAddr      string // server address
	Port        uint64 // port
	Scheme      string // http or https
	ContextPath string // path
	NamespaceID string // namespace id
	// if the parameter is not empty, replace the same fields as ClientConfig and ServerConfig above
	clientConfig  *constant.ClientConfig
	serverConfigs []constant.ServerConfig

	Group  string // group, example: dev, prod, test
	DataID string // config file id
	Format string // configuration file type: json,yaml,toml
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

	// create clientConfig
	if params.clientConfig == nil {
		params.clientConfig = &constant.ClientConfig{
			NamespaceId:         params.NamespaceID,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              os.TempDir() + "/nacos/log",
			CacheDir:            os.TempDir() + "/nacos/cache",
		}
	}

	// create serverConfig
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

// Init get configuration from nacos and parse to struct
func Init(obj interface{}, params *Params, opts ...Option) error {
	err := params.valid()
	if err != nil {
		return err
	}

	setParams(params, opts...)

	// create a dynamic configuration client
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  params.clientConfig,
			ServerConfigs: params.serverConfigs,
		},
	)
	if err != nil {
		return err
	}

	// read config content
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: params.DataID,
		Group:  params.Group,
	})
	if err != nil {
		return err
	}

	// parse config
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

// NewNamingClient create a service registration and discovery of nacos client
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
