// Package consulcli is connecting to the consul service client.
package consulcli

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// Init connecting to the consul service
// Note: If the WithConfig(*api.Config) parameter is set, the addr parameter is ignored!
func Init(addr string, opts ...Option) (*api.Client, error) {
	o := defaultOptions()
	o.apply(opts...)

	if o.config != nil {
		return api.NewClient(o.config)
	}

	if addr == "" {
		return nil, fmt.Errorf("consul address cannot be empty")
	}

	return api.NewClient(&api.Config{
		Address:    addr,
		Scheme:     o.scheme,
		WaitTime:   o.waitTime,
		Datacenter: o.datacenter,
	})
}
