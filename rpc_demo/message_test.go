package rpc_demo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeMsg(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "with meta",
			req: &Request{
				MessId:      123,
				Version:     12,
				Compresser:  25,
				Serializer:  17,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "b",
					"shadow":   "true",
				},
				Arg: []byte("hello, world"),
			},
		},
		{
			name: "no meta",
			req: &Request{
				MessId:      123,
				Version:     12,
				Compresser:  25,
				Serializer:  17,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Arg:         []byte("hello, world"),
			},
		},
		{
			name: "empty value",
			req: &Request{
				MessId:      123,
				Version:     12,
				Compresser:  25,
				Serializer:  17,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "",
					"shadow":   "true",
				},
				Arg: []byte("hello, world"),
			},
		},
	}
	for _, tc := range testCases {
		// 这里测试我们利用 encode/decode 过程相反的特性
		t.Run(tc.name, func(t *testing.T) {
			CalRequestHead(tc.req)
			tc.req.BodyLength = uint32(len(tc.req.Arg))
			bs, _ := EncodeReq(tc.req)
			req, _ := DecodeReq(bs)
			assert.Equal(t, tc.req, req)
		})
	}

}
