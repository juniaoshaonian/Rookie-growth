package main

import (
	"fmt"
	"im_decoder/frame"
	"im_decoder/Protocol"
	"log"
	"net"
)
func handlePackage(goimBody []byte)*Protocol.Goim{
	 p := &Protocol.Goim{}
	err := p.Decode(goimBody)
	if err != nil {
		fmt.Println("goim decode error")
		return nil
	}
	return p

}

func handleConn(c net.Conn){
	defer c.Close()
	goimCoder := frame.NewGoimCoder()
	for {
		goimdata,err := goimCoder.Decode(c)
		if err != nil {
			fmt.Println("handleConnERROR")
			return
		}
		goim := handlePackage(goimdata)
		fmt.Println("goim is:",goim)

	}

}
func main() {
	listen,err := net.Listen("tcp","localhost:8000")
	if err != nil {
		log.Println(err)
	}
	defer listen.Close()
	for {
		conn,err := listen.Accept()
		if err != nil {
			log.Println(conn.RemoteAddr().String(),"发生了错误",err)
			continue
		}
		go handleConn(conn)
	}
}