package roundrobin

import (
	"errors"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"rpc/loadbalance"
	"sync"
)

type WeightPickBuilder struct {
}

func (p *WeightPickBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*Conn, 0, len(info.ReadySCs))
	for subconn, val := range info.ReadySCs {
		weight := val.Address.Attributes.Value("weight").(uint32)
		conn := &Conn{
			address:         val.Address,
			SubConn:         subconn,
			weight:          weight,
			efficientWeight: weight,
			currentWeight:   weight,
		}
		conns = append(conns, conn)
	}
	return &WeightPick{
		conns:  conns,
		filter: loadbalance.GroupFilter,
	}

}

type WeightPick struct {
	wg     sync.RWMutex
	conns  []*Conn
	filter loadbalance.Filter
}

func (p *WeightPick) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	res := balancer.PickResult{}
	if len(p.conns) == 0 {
		return res, errors.New("未找到可用节点")
	}
	p.wg.Lock()
	defer p.wg.Unlock()
	var candidate *Conn
	conns := make([]*Conn, 0, len(p.conns))
	for _, cc := range p.conns {
		if p.filter(info.Ctx, cc.address) {
			conns = append(conns, cc)
		}
	}

	totalWeight := uint32(0)
	for _, cc := range conns {
		totalWeight += cc.currentWeight
		cc.currentWeight += cc.efficientWeight
		if candidate == nil || cc.currentWeight > candidate.currentWeight {
			candidate = cc
		}
	}

	candidate.currentWeight = candidate.currentWeight - totalWeight

	return balancer.PickResult{
		SubConn: candidate.SubConn,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Conn struct {
	balancer.SubConn
	address         resolver.Address
	weight          uint32
	efficientWeight uint32
	currentWeight   uint32
}
