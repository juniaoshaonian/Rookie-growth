package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	R *http.Request
	W http.ResponseWriter
}
//将body序列化为对象
func (c *Context)ReadJSON(req interface{})error{
	r:=c.R
	body,err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body,req)
	if err != nil {
		return err
	}
	return nil
}
func (c *Context)WriteJSON(code int,resp interface{})error{
	c.W.WriteHeader(code)
	data,err := json.Marshal(resp)
	if err!= nil {
		return err
	}
	_,err = c.W.Write(data)
	return err
}
func (c *Context)Badresp(resp interface{})error{
	return c.WriteJSON(http.StatusInternalServerError,resp)
}
func NewContext(r *http.Request,w http.ResponseWriter)*Context{
	return &Context{
		R: r,
		W: w,
	}
}