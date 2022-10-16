package accesslog

import (
	webframe_22 "webframe"
)

type MiddlewareBuilder struct {
	logFunc func(accesssLog string)
}

func (m *MiddlewareBuilder) LogFunc(logfunc func(accessLog string)) *MiddlewareBuilder {
	m.logFunc = logfunc
	return m
}

type acccessLog struct {
	Host   string
	Route  string
	Method string `json:"http_method"`
	Path   string
}

func (m *MiddlewareBuilder) Build() webframe_22.Middleware {
	return func(next webframe_22.HandleFunc) webframe_22.HandleFunc {
		return func(c *webframe_22.Context) {
			defer func() {
				l := acccessLog{
					Host: c.Req.Host,
				}
			}()
		}

	}
}
