package nacoscli

import (
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

type options struct {
	clientConfig  *constant.ClientConfig
	serverConfigs []constant.ServerConfig
}

func defaultOptions() *options {
	return &options{
		clientConfig:  nil,
		serverConfigs: nil,
	}
}

// Option set the jwt options.
type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithClientConfig set nacos client config
func WithClientConfig(clientConfig *constant.ClientConfig) Option {
	return func(o *options) {
		o.clientConfig = clientConfig
	}
}

// WithServerConfigs set nacos server config
func WithServerConfigs(serverConfigs []constant.ServerConfig) Option {
	return func(o *options) {
		o.serverConfigs = serverConfigs
	}
}
