package resolver

import (
	"github.com/xeitf/lamp"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	serviceName string
	closeWatch  func()
	clientConn  resolver.ClientConn
}

// NewResolver
func NewResolver(serviceName string, cc resolver.ClientConn) (r *Resolver) {
	return &Resolver{serviceName: serviceName, clientConn: cc}
}

// ResolveNow implements resolver.Resolver.
func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {

}

// Watch
func (r *Resolver) Watch() (err error) {
	r.closeWatch, err = lamp.Watch(r.serviceName, "grpc", r.Update)
	return
}

// Update
func (r *Resolver) Update(addrs []string, closed bool) {
	var grpcAddrs []resolver.Address
	for _, addr := range addrs {
		grpcAddrs = append(grpcAddrs, resolver.Address{Addr: addr})
	}
	r.clientConn.UpdateState(resolver.State{Addresses: grpcAddrs})
}

// Close implements resolver.Resolver.
func (r *Resolver) Close() {
	if r.closeWatch != nil {
		r.closeWatch()
	}
}

type Builder struct {
}

// NewBuilder
func NewBuilder() (b *Builder) {
	return &Builder{}
}

// Build implements resolver.Builder.
func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (r resolver.Resolver, err error) {
	nr := NewResolver(target.Endpoint(), cc)
	if err = nr.Watch(); err != nil {
		return nil, err
	}
	return nr, nil
}

// Scheme implements resolver.Builder.
func (b *Builder) Scheme() string {
	return "lamp"
}

// Register
func Register() {
	resolver.Register(NewBuilder())
}
