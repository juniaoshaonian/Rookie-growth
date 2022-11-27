package rpc_demo

import "context"

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

type Request struct {
	// 请求头长度
	HeadLength uint32
	// 请求体长度
	BodyLength uint32
	MessId     uint32
	Compresser byte
	Serializer byte
	Version    byte
	Meta       map[string]string
	// 服务名
	ServiceName string
	// 方法名
	MethodName string
	// 参数
	Arg []byte
}

type Response struct {
	HeadLength uint32
	BodyLength uint32

	MessageId uint32

	Version    byte
	Compresser byte
	Serializer byte

	Error []byte
	// 你要区分业务 error 还是非业务 error
	// BizError []byte // 代表的是业务返回的 error

	Data []byte
}

func CalRequestHead(request *Request) {
	ans := 0
	base := 15
	ans += base
	ans = ans + len(request.ServiceName) + 1
	ans = ans + len(request.MethodName) + 1
	for key, val := range request.Meta {
		ans = ans + len(key) + 1 + len(val) + 1
	}
	request.HeadLength = uint32(ans)
}
