package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"rpc_demo2/register"
	"sync"
)

var typesMap = map[mvccpb.Event_EventType]register.EventType{
	mvccpb.PUT:    register.EventTypeAdd,
	mvccpb.DELETE: register.EventTypeDelete,
}

type EtcdRegister struct {
	mu          sync.RWMutex
	client      *clientv3.Client
	cancelFuncs []context.CancelFunc
	sess        concurrency.Session
}

func (r *EtcdRegister) Register(ctx context.Context, ins register.ServiceInstance) error {
	serviceKey := fmt.Sprintf("/mirco/%s/%s", ins.ServiceName, ins.Addr)
	insdata, err := json.Marshal(ins)
	if err != nil {
		return err
	}
	_, err = r.client.Put(ctx, serviceKey, string(insdata), clientv3.WithLease(r.sess.Lease()))
	if err != nil {
		return err
	}

	return nil
}

func (r *EtcdRegister) UnRegiter(ctx context.Context, ins register.ServiceInstance) error {
	instanceKey := fmt.Sprintf("/micro/%s/%s", ins.ServiceName, ins.Addr)
	_, err := r.client.Delete(ctx, instanceKey)
	return err
}

func (r *EtcdRegister) ListService(ctx context.Context, serviceName string) ([]*register.ServiceInstance, error) {
	serviceKey := fmt.Sprintf("/mirco/%s", serviceName)
	resp, err := r.client.Get(ctx, serviceKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	res := make([]*register.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		ins := &register.ServiceInstance{}
		err = json.Unmarshal(kv.Value, ins)
		if err != nil {
			continue
		}
		res = append(res, ins)
	}
	return res, nil
}

func (r *EtcdRegister) SubScribe(serviceName string) (<-chan register.Event, error) {
	serviceKey := fmt.Sprintf("/mirco/%s", serviceName)

	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFuncs = append(r.cancelFuncs, cancel)
	watchCh := r.client.Watch(context.Background(), serviceKey, clientv3.WithPrefix())
	eventCh := make(chan register.Event)
	// 开启一个协程监听 watchCh,怎样关闭这个协程
	go func() {
		for watcher := range watchCh {
			if ctx.Err() != nil {
				return
			}
			if watcher.Err() != nil {
				continue
			}
			for _, event := range watcher.Events {
				ins := &register.ServiceInstance{}
				err := json.Unmarshal(event.Kv.Value, ins)
				if err != nil {
					select {
					case eventCh <- register.Event{}:
					case <-ctx.Done():
						return
					}
					continue
				}
				select {
				case eventCh <- register.Event{
					Type: typesMap[event.Type],
				}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return eventCh, nil
}

func (r *EtcdRegister) Close() error {
	r.mu.RLock()
	cancelFuncs := r.cancelFuncs
	r.mu.RUnlock()

	for _, cancel := range cancelFuncs {
		cancel()
	}
	return nil
}
