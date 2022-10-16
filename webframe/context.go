package webframe

import (
	"net/http"
)

type Context struct {
	Req           *http.Request
	Resp          http.ResponseWriter
	RespsonseDate []byte
	ResponseCode  int
	MatchRouter   string
	pathParams    map[string]string
	UserValues    map[string]any
}

func (ctx *Context) QueryValue(key string) (string, error) {
	params := ctx.Req.URL.Query()
	val, ok := params[key]
	if !ok || len(val) == 0 {
		return "", nil
	}
	return val[0], nil
}
