package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

type PackageBody []byte
type Coder interface {
	Encode(writer io.Writer,body PackageBody)error
	Decode(reader io.Reader)(PackageBody,error)
}

type GoimCoder struct {

}
func  NewGoimCoder()Coder{
	return &GoimCoder{}
}
func (c *GoimCoder)Encode(writer io.Writer,body PackageBody)error {
	var length int32 = int32(len(body)) + 4
	err := binary.Write(writer,binary.BigEndian,&length)
	if err != nil {
		return err
	}
	n,err := writer.Write(body)
	if err != nil {
		return err
	}
	if n != len(body) {
		return errors.New("short write error")
	}
	return nil
}

func (c *GoimCoder)Decode(r io.Reader)(PackageBody,error){
	var length int32
	err:=binary.Read(r,binary.BigEndian,&length)
	if err != nil {
		return nil,err
	}
	buf := make([]byte,length - 4)
	n,err := io.ReadFull(r,buf)
	if err != nil {
		return nil,err
	}
	if n != int(length -4){
		return nil,errors.New("short read")
	}
	return PackageBody(buf),nil
}
