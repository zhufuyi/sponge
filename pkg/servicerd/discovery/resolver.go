// Package discovery is service discovery library, supports etcd, consul and nacos.
package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/go-dev-frame/sponge/pkg/servicerd/registry"
)

type discoveryResolver struct {
	w  registry.Watcher
	cc resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	insecure         bool
	debugLogDisabled bool
}

func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			fmt.Printf("[resolver] Failed to watch discovery endpoint: %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*registry.ServiceInstance) {
	addrs := make([]resolver.Address, 0)
	endpoints := make(map[string]struct{})
	for _, in := range ins {
		endpoint, err := parseEndpoint(in.Endpoints, "grpc", !r.insecure)
		if err != nil {
			//fmt.Printf("[resolver] Failed to parse discovery endpoint: %v\n", err)
			continue
		}
		if endpoint == "" {
			continue
		}
		// filter redundant endpoints
		if _, ok := endpoints[endpoint]; ok {
			continue
		}
		endpoints[endpoint] = struct{}{}
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.Metadata),
			Addr:       endpoint,
		}
		addr.Attributes = addr.Attributes.WithValue("rawServiceInstance", in)
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		//fmt.Printf("[resolver] Zero endpoint found,refused to write, instances: %v\n", ins)
		return
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		fmt.Printf("[resolver] failed to update state: %v\n", err)
	}

	if !r.debugLogDisabled {
		b, _ := json.Marshal(ins)
		fmt.Printf("[resolver] update instances: %s\n", b)
	}
}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.w.Stop()
	if err != nil {
		fmt.Printf("[resolver] failed to watch top: %v\n", err)
	}
}

func (r *discoveryResolver) ResolveNow(_ resolver.ResolveNowOptions) {}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}

// parseEndpoint parses an Endpoint URL.
func parseEndpoint(endpoints []string, scheme string, isSecure bool) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if u.Scheme == scheme && IsSecure(u) == isSecure {
			return u.Host, nil
		}
	}
	return "", nil
}

// IsSecure parses isSecure for Endpoint URL.
func IsSecure(u *url.URL) bool {
	ok, err := strconv.ParseBool(u.Query().Get("isSecure"))
	if err != nil {
		return false
	}
	return ok
}
