package rpc_demo

import (
	"context"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"time"
)

type Client struct {
	pool pool.Pool
}

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	co, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	conn := co.(net.Conn)

	data, err := EncodeReq(req)
	if err != nil {
		return nil, err
	}
	n, err := conn.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, errors.New("tcp 未写入全部数据")
	}
	respdata, err := ReadMsg(conn)
	if err != nil {
		return nil, err
	}
	return DecodeResp(respdata), nil
}

func NewClient(addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,
		MaxCap:     100,
		MaxIdle:    50,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		},
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
		IdleTimeout: time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		pool: p,
	}, nil
}
