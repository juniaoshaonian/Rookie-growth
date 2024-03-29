package rpc

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"rpc/compress"
	"rpc/message"
	"rpc/serialize"
	"rpc/serialize/json"
)

type Server struct {
	services    map[string]*reflectionStub
	serializers []serialize.Serializer
	compressers []compress.Compresser
}

func (s *Server) Start(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept connection got error: %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		bytes, err := ReadMsg(conn)
		if err != nil {
			return
		}
		req := message.DecodeReq(bytes)
		resp, er := s.Invoke(context.Background(), req)
		if er != nil {
			resp = &message.Response{}
			// 服务器本身出错了
			resp.Error = []byte(fmt.Errorf("rpc-server: 服务器异常 %w", er).Error())
			// 计算一下长度
			resp.SetHeadLength()
		}
		encode := message.EncodeResp(resp)
		_, er = conn.Write(encode)
		if er != nil {
			fmt.Printf("sending response failed: %v", er)
		}
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	resp := &message.Response{}
	service, ok := s.services[req.ServiceName]
	if !ok {
		return resp, fmt.Errorf("server: 未找到服务, 服务名 %s", req.ServiceName)
	}
	return service.invoke(ctx, req)
}

func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Code()] = serializer
}
func (s *Server) RegisterCompresser(compresser compress.Compresser) {
	s.compressers[compresser.Code()] = compresser
}

func (s *Server) RegisterService(service Service) {
	s.services[service.ServiceName()] = &reflectionStub{
		s:           service,
		serializers: s.serializers,
		value:       reflect.ValueOf(service),
		compressers: s.compressers,
	}
}

func NewServer() *Server {
	res := &Server{
		services: make(map[string]*reflectionStub, 4),
		// 一个字节，最多有 256 个实现，直接做成一个简单的 bit array 的东西
		serializers: make([]serialize.Serializer, 256),
		compressers: make([]compress.Compresser, 256),
	}
	// 注册最基本的序列化协议
	res.RegisterSerializer(json.Serializer{})
	return res
}

type reflectionStub struct {
	s           Service
	value       reflect.Value
	serializers []serialize.Serializer
	compressers []compress.Compresser
}

func (r *reflectionStub) invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	var err error
	method := r.value.MethodByName(req.Method)
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	c := r.compressers[req.Compresser]
	req.Data, err = c.Decompress(req.Data)
	if err != nil {
		return nil, err
	}
	s := r.serializers[req.Serializer]
	err = s.Decode(req.Data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})
	respData, err := s.Encode(res[0].Interface())
	if err != nil {
		// 服务器本身的错误
		return nil, err
	}
	respData, err = c.Compress(respData)
	if err != nil {
		return nil, err
	}
	resp := &message.Response{
		BodyLength: uint32(len(respData)),
		MessageId:  req.MessageId,
		Compresser: req.Compresser,
		// 理论上来说，这里可以换一种序列化协议，但是没必要暴露这种功能给用户
		Serializer: req.Serializer,
		Data:       respData,
	}
	if !res[1].IsZero() {
		resp.Error = []byte(res[1].Interface().(error).Error())
	}
	resp.SetHeadLength()
	return resp, nil
}
