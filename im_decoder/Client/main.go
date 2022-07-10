package main

import (
	"fmt"
	"im_decoder/frame"
	"im_decoder/Protocol"
	"log"
	"net"
	"os"
	"strconv"
)
func sendData(conn net.Conn){
	for i:=0;i<5;i++{
		 p := &Protocol.Goim{}
		 p.SequenceId = 10
		 p.Protocolversion = 2
		 p.Operation = 99
		p.Body = []byte(strconv.Itoa(i))
		word ,err :=  p.Encode()
		if err != nil {
			log.Println(err)
		}
		f := frame.NewGoimCoder()
		err = f.Encode(conn,word)
		if err != nil {
			log.Println(err)
		}
	}
}

func main(){
	server := "localhost:8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	sendData(conn)
}