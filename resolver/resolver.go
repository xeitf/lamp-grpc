package resolver

import (
	"github.com/xeitf/lamp"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
}

// ResolveNow implements resolver.Resolver.
func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {

}

// Close implements resolver.Resolver.
func (r *Resolver) Close() {

}

type Builder struct {
}

// NewBuilder
func NewBuilder() (b *Builder) {
	return &Builder{}
}

// Scheme implements resolver.Builder.
func (b *Builder) Scheme() string {
	return "lamp"
}

// Build implements resolver.Builder.
func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (r resolver.Resolver, err error) {
	lamp.Watch(target.Endpoint(), "grpc", func(addrs []string, closed bool) {
		var grpcAddrs []resolver.Address
		for _, addr := range addrs {
			grpcAddrs = append(grpcAddrs, resolver.Address{Addr: addr})
		}
		cc.UpdateState(resolver.State{Addresses: grpcAddrs})
	})
	return &Resolver{}, nil
}

// Init
func Init(cfg string) (close func() error, err error) {
	close, err = lamp.Init(cfg)
	if err == nil {
		resolver.Register(NewBuilder())
	}
	return
}
