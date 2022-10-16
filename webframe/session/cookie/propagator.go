package cookie

import (
	"net/http"
)

type PropagatorOpt func(p *Propagator)
type CookieOpt func(cookie *http.Cookie)

func WithCookieopt(c CookieOpt) PropagatorOpt {
	return func(p *Propagator) {
		p.cookieOpt = c
	}
}

func NewPropagator(cookieName string, opts ...PropagatorOpt) *Propagator {
	p := &Propagator{
		cookieName: cookieName,
		cookieOpt: func(cookie *http.Cookie) {

		},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type Propagator struct {
	cookieName string
	cookieOpt  CookieOpt
}

func (p *Propagator) Inject(id string, resp http.ResponseWriter) error {
	// 将session id 写入响应
	c := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOpt(c)
	http.SetCookie(resp, c)
	return nil

}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	val, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return val.Value, nil

}

func (p *Propagator) Remove(resp http.ResponseWriter) error {
	http.SetCookie(resp, &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	})
	return nil
}
