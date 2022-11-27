package rpc_demo

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	Services map[string]Service
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Register(service Service) {
	if s.Services == nil {
		s.Services = make(map[string]Service)
	}
	s.Services[service.Name()] = service
}
func (s *Server) Invoke(ctx context.Context, req *Request) (*Response, error) {
	resp := &Response{
		MessageId:  req.MessId,
		Version:    req.Version,
		Compresser: req.Compresser,
		Serializer: req.Serializer,
	}
	service, ok := s.Services[req.ServiceName]
	if !ok {
		resp.Error = []byte(errors.New("未找到相关方法").Error())
		resp.SetHeadLength()
		return resp, errors.New("未找到相关方法")
	}
	serviceVal := reflect.ValueOf(service)
	// 找到返回值的类型

	method := serviceVal.MethodByName(req.MethodName)
	//初始化参数
	arg := reflect.New(method.Type().In(1))
	// 反序列化参数给arg赋值
	err := json.Unmarshal(req.Arg, arg.Interface())
	if err != nil {
		resp.Error = []byte(err.Error())
		resp.SetHeadLength()
		return resp, err
	}
	vals := method.Call([]reflect.Value{reflect.ValueOf(ctx), arg.Elem()})
	if len(vals) > 1 && !vals[1].IsZero() {
		resp.Error = []byte(vals[1].Interface().(error).Error())
		resp.SetHeadLength()
		return resp, vals[1].Interface().(error)
	}
	data, err := json.Marshal(vals[0].Interface())
	if err != nil {
		resp.Error = []byte(err.Error())
		resp.SetHeadLength()
		return resp, err
	}
	resp.BodyLength = uint32(len(data))
	resp.Data = data
	resp.SetHeadLength()
	return resp, nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go func() error {
			err := s.HandleConn(conn)
			if err != nil {
				conn.Close()
				return err
			}
			return nil
		}()

	}
}

func (s *Server) HandleConn(conn net.Conn) error {
	for {
		data, err := ReadMsg(conn)
		if err != nil {
			return err
		}
		req, err := DecodeReq(data)
		if err != nil {
			return err
		}
		resp, err := s.Invoke(context.Background(), req)
		if err != nil {
			return err
		}
		encodeData := EncodeResp(resp)
		_, err = conn.Write(encodeData)
		if err != nil {
			return err
		}
		return nil
	}

}
