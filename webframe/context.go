package webframe

import "net/http"

type Context struct {
	Req           *http.Request
	Resp          http.ResponseWriter
	RespsonseDate []byte
	ResponseCode  int
	MatchRouter   string
}
