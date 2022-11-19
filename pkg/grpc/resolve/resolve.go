package resolve

import (
	"fmt"
	"net/url"
	"sync"

	"google.golang.org/grpc/resolver"
)

var mutex = &sync.Mutex{}

// Register address and serviceName
func Register(scheme string, serviceName string, address []string) string {
	mutex.Lock()
	defer mutex.Unlock()

	endpoint := fmt.Sprintf("%s:///%s", scheme, serviceName)
	u, _ := url.Parse(endpoint)

	resolver.Register(&ResolverBuilder{
		scheme:      scheme,
		serviceName: serviceName,
		addrs:       address,
		path:        u.Path,
	})

	return endpoint
}

// ResolverBuilder resolver struct
type ResolverBuilder struct {
	scheme      string
	serviceName string
	addrs       []string
	path        string
}

// Build resolver
func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	blr := &blResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			r.path: r.addrs,
		},
	}
	blr.start()
	return blr, nil
}

// Scheme get scheme
func (r *ResolverBuilder) Scheme() string {
	return r.scheme
}

type blResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (b *blResolver) start() {
	addrStrs := b.addrsStore[b.target.URL.Path]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	_ = b.cc.UpdateState(resolver.State{Addresses: addrs})
}

// ResolveNow Resolve now
func (*blResolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close resolver
func (*blResolver) Close() {}
