package consulcli

import (
	"github.com/hashicorp/consul/api"
)

// Init connecting to the consul service
func Init(addr string, opts ...Option) (*api.Client, error) {
	o := defaultOptions()
	o.apply(opts...)

	return api.NewClient(&api.Config{
		Address:    addr,
		Scheme:     o.scheme,
		WaitTime:   o.waitTime,
		Datacenter: o.datacenter,
	})
}
