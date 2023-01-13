package rpc

import (
	"context"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"rpc/registry"
)

type grpcResolverBuilder struct {
	r registry.Registry
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	res := &grpcResolver{
		target:  target,
		cc:      cc,
		r:       g.r,
		closech: make(chan struct{}),
	}
	res.resolve()
	return res, res.watch()
}

func (g *grpcResolverBuilder) Scheme() string {
	return "register"
}

func NewResolverBuilder(r registry.Registry) resolver.Builder {
	return &grpcResolverBuilder{
		r: r,
	}
}

type grpcResolver struct {
	r       registry.Registry
	target  resolver.Target
	cc      resolver.ClientConn
	closech chan struct{}
}

//立刻執行服務發現
func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) Close() {
	g.closech <- struct{}{}
}

func (g *grpcResolver) resolve() {
	r := g.r
	instances, err := r.ListService(context.Background(), g.target.Endpoint)
	if err != nil {
		g.cc.ReportError(err)
		return
	}

	addresses := make([]resolver.Address, 0, len(instances))
	for _, ins := range instances {
		addresses = append(addresses, resolver.Address{
			Addr:       ins.Address,
			Attributes: attributes.New("weight", ins.Weight).WithValue("group", ins.Group),
		})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
	if err != nil {
		g.cc.ReportError(err)
	}
}

func (g *grpcResolver) watch() error {
	eventCh, err := g.r.Subscribe(g.target.Endpoint)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-eventCh:
				g.resolve()
			case <-g.closech:
				close(g.closech)
				return
			}
		}
	}()
	return nil
}
