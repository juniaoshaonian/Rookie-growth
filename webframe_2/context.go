package webframe_2

import "net/http"

type Context struct {
	Req      *http.Request
	Resp     http.ResponseWriter
	paramMap map[string]string
}
