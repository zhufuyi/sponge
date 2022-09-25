package loadbalance

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/resolver"
)

var mux = &sync.Mutex{}

// Register 注册地址到resolver map
func Register(schemeStr string, serviceNameStr string, address []string) string {
	mux.Lock()
	defer mux.Unlock()

	resolver.Register(&ResolverBuilder{
		SchemeVal:   schemeStr,
		ServiceName: serviceNameStr,
		Addrs:       address,
	})

	return fmt.Sprintf("%s:///%s", schemeStr, serviceNameStr)
}

// ResolverBuilder 解析生成器
type ResolverBuilder struct {
	SchemeVal   string // SchemeVal作为唯一标致，重复会被覆盖addrs
	ServiceName string
	Addrs       []string
}

// Build 创建生成器
func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	blr := &blResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			r.ServiceName: r.Addrs,
		},
	}
	blr.start()
	return blr, nil
}

// Scheme 设置Scheme值
func (r *ResolverBuilder) Scheme() string {
	return r.SchemeVal
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

// ResolveNow 当前解析生成器
func (*blResolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close 关闭
func (*blResolver) Close() {}
