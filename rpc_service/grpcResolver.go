package rpc_service

import (
	"context"
	"google.golang.org/grpc/resolver"
)

type grpcResolverBuilder struct {
	r Registry
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	subch, err := g.r.Subscribe(target.Endpoint)
	if err != nil {
		return nil, err
	}
	res := &grpcResolver{
		r:      g.r,
		target: target,
		cc:     cc,
	}
	go func() {
		for {
			select {
			case <-subch:
				res.resolve()
			}
		}
	}()
	return res, nil
}
func NewResolverBuilder() resolver.Builder {
	return &grpcResolverBuilder{}
}

func (g *grpcResolverBuilder) Scheme() string {
	//TODO implement me
	panic("implement me")
}

type grpcResolver struct {
	r      Registry
	target resolver.Target
	cc     resolver.ClientConn
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) Close() {
	//TODO implement me
	panic("implement me")
}

func (g *grpcResolver) resolve() {
	r := g.r
	inss, err := r.ListService(context.Background(), g.target.Endpoint)
	if err != nil {

	}
	addresses := make([]resolver.Address, 0, len(inss))
	for _, ins := range inss {
		addresses = append(addresses, resolver.Address{
			Addr: ins.Addr,
		})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
	if err != nil {

	}
}
