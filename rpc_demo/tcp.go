package rpc_demo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

const (
	maxbytes     int = 8
	splitter         = '\n'
	pairSplitter     = '\r'
)

func EncodeReq(req *Request) ([]byte, error) {
	ds := make([]byte, req.HeadLength+req.BodyLength)
	cur := 0
	binary.BigEndian.PutUint32(ds[cur:cur+4], req.HeadLength)
	cur = cur + 4
	binary.BigEndian.PutUint32(ds[cur:cur+4], req.BodyLength)
	cur = cur + 4
	binary.BigEndian.PutUint32(ds[cur:cur+4], req.MessId)
	cur = cur + 4
	ds[cur] = req.Version
	cur++
	ds[cur] = req.Compresser
	cur++
	ds[cur] = req.Serializer
	cur++
	copy(ds[cur:cur+len(req.ServiceName)], req.ServiceName)
	cur += len(req.ServiceName)
	ds[cur] = splitter
	cur++
	copy(ds[cur:cur+len(req.MethodName)], req.MethodName)
	cur += len(req.MethodName)
	ds[cur] = splitter
	cur++
	for key, val := range req.Meta {
		copy(ds[cur:cur+len(key)], key)
		cur += len(key)
		ds[cur] = pairSplitter
		cur++
		copy(ds[cur:cur+len(val)], val)
		cur += len(val)
		ds[cur] = splitter
		cur++
	}
	copy(ds[cur:], req.Arg)
	return ds, nil
}

func DecodeReq(data []byte) (*Request, error) {
	req := &Request{}
	req.HeadLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.MessId = binary.BigEndian.Uint32(data[8:12])
	req.Version = data[12]
	req.Compresser = data[13]
	req.Serializer = data[14]
	head := data[15:req.HeadLength]
	index := bytes.IndexByte(head, splitter)
	req.ServiceName = string(head[:index])
	head = head[index+1:]
	index = bytes.IndexByte(head, splitter)
	req.MethodName = string(head[:index])
	head = head[index+1:]

	for len(head) > 0 {
		if req.Meta == nil {
			req.Meta = make(map[string]string)
		}
		index := bytes.IndexByte(head, splitter)
		pairindex := bytes.IndexByte(head, pairSplitter)
		key := head[:pairindex]
		val := head[pairindex+1 : index]

		req.Meta[string(key)] = string(val)
		if index+1 >= len(head) {
			break
		}
		head = head[index+1:]
	}
	req.Arg = data[req.HeadLength:]
	return req, nil
}

func (resp *Response) SetHeadLength() {
	resp.HeadLength = uint32(15 + len(resp.Error))
}

// 这里处理 Resp 我直接复制粘贴，是因为我觉得复制粘贴会使可读性更高

func EncodeResp(resp *Response) []byte {
	bs := make([]byte, resp.HeadLength+resp.BodyLength)

	cur := bs
	// 1. 写入 HeadLength，四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.HeadLength)
	cur = cur[4:]
	// 2. 写入 BodyLength 四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]

	// 3. 写入 message id, 四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.MessageId)
	cur = cur[4:]

	// 4. 写入 version，因为本身就是一个字节，所以不用进行编码了
	cur[0] = resp.Version
	cur = cur[1:]

	// 5. 写入压缩算法
	cur[0] = resp.Compresser
	cur = cur[1:]

	// 6. 写入序列化协议
	cur[0] = resp.Serializer
	cur = cur[1:]
	// 7. 写入 error
	copy(cur, resp.Error)
	cur = cur[len(resp.Error):]

	// 剩下的数据
	copy(cur, resp.Data)
	return bs
}

// DecodeResp 解析 Response
func DecodeResp(bs []byte) *Response {
	resp := &Response{}
	// 按照 EncodeReq 写下来
	// 1. 读取 HeadLength
	resp.HeadLength = binary.BigEndian.Uint32(bs[:4])
	// 2. 读取 BodyLength
	resp.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	// 3. 读取 message id
	resp.MessageId = binary.BigEndian.Uint32(bs[8:12])
	// 4. 读取 Version
	resp.Version = bs[12]
	// 5. 读取压缩算法
	resp.Compresser = bs[13]
	// 6. 读取序列化协议
	resp.Serializer = bs[14]
	// 7. error 信息
	resp.Error = bs[15:resp.HeadLength]

	// 剩下的就是数据了
	resp.Data = bs[resp.HeadLength:]
	return resp
}

func ReadMsg(conn net.Conn) ([]byte, error) {
	length := make([]byte, maxbytes)
	n, err := conn.Read(length)
	if err != nil {
		return nil, err
	}
	if n != maxbytes {
		return nil, errors.New("长度字段没有读全")
	}

	Headlength := binary.BigEndian.Uint32(length[:4])
	Bodylength := binary.BigEndian.Uint32(length[4:])
	data := make([]byte, Headlength+Bodylength)
	copy(data[:8], length)
	_, err = conn.Read(data[8:])
	if err != nil {
		return nil, err
	}
	return data, nil

}
