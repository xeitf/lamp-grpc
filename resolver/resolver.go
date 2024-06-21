package resolver

import (
	"github.com/xeitf/lamp"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	lc          *lamp.Client
	serviceName string
	closeWatch  func()
	clientConn  resolver.ClientConn
}

// Watch
func (r *Resolver) Watch() (err error) {
	r.closeWatch, err = r.lc.Watch(r.serviceName, "grpc", r.update)
	return
}

// update
func (r *Resolver) update(addrs []lamp.Address, closed bool) {
	var grpcAddrs []resolver.Address
	for _, addr := range addrs {
		grpcAddrs = append(grpcAddrs, resolver.Address{Addr: addr.Addr})
	}
	r.clientConn.UpdateState(resolver.State{Addresses: grpcAddrs})
}

// ResolveNow implements resolver.Resolver.
func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {

}

// Close implements resolver.Resolver.
func (r *Resolver) Close() {
	if r.closeWatch != nil {
		r.closeWatch()
	}
}

type Builder struct {
	lc *lamp.Client
}

// NewBuilder
func NewBuilder(lc *lamp.Client) (b *Builder) {
	return &Builder{lc: lc}
}

// Scheme implements resolver.Builder.
func (b *Builder) Scheme() string {
	return "lamp"
}

// Build implements resolver.Builder.
func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (r resolver.Resolver, err error) {
	nr := &Resolver{
		lc:          b.lc,
		clientConn:  cc,
		serviceName: target.Endpoint(),
	}
	if err = nr.Watch(); err != nil {
		return nil, err
	}
	return nr, nil
}

// Register
func Register(lc *lamp.Client) {
	resolver.Register(NewBuilder(lc))
}
