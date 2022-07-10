package Protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Protocol interface {
	Decode([]byte) error
	Encode()([]byte,error)
}

const (
	HeaderLengthBytes = 2
	ProtocolVersionBytes = 2
	OperationBytes = 4
	SequenceIDBytes = 4
	HeaderLength =  HeaderLengthBytes + ProtocolVersionBytes + OperationBytes + SequenceIDBytes
)
type Goim struct {

	Protocolversion int
	Operation int
	SequenceId int
	Body []byte
}
func ByteToInt16(n []byte)int{
	bytesbuffer := bytes.NewBuffer(n)
	var x int16
	binary.Read(bytesbuffer,binary.BigEndian,&x)
	return int(x)
}
func ByteToInt32(n []byte)int{
	bytesbuffer := bytes.NewBuffer(n)
	var x int32
	binary.Read(bytesbuffer,binary.BigEndian,&x)
	return int(x)
}
func (g *Goim)Decode(data []byte)error{
	site := HeaderLengthBytes
	headerLength := ByteToInt16(data[:site])
	if headerLength != HeaderLength {
		return errors.New("headerLength error")
	}
	g.Protocolversion = ByteToInt16(data[site:site+ProtocolVersionBytes])
	site +=ProtocolVersionBytes
	g.Operation = ByteToInt32(data[site:site+OperationBytes])
	site += OperationBytes
	g.SequenceId = ByteToInt32(data[site:site+SequenceIDBytes])
	site += SequenceIDBytes
	g.Body = data[site:]

	return nil
}
func Int32ToBytes(n int)[]byte{
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer,binary.BigEndian,x)
	return bytesBuffer.Bytes()
}
func Int16ToBytes(n int)[]byte {
	x := int16(n)
	bytesBuffers := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffers,binary.BigEndian,x)
	return bytesBuffers.Bytes()
}

func (g *Goim)Encode()([]byte,error){
  var data []byte
  data = append(data,Int16ToBytes(HeaderLength)...)
  data = append(data,Int16ToBytes(g.Protocolversion)...)
  data = append(data,Int32ToBytes(g.Operation)...)
  data = append(data,Int32ToBytes(g.SequenceId)...)
  data = append(data,g.Body...)
  return data,nil

}