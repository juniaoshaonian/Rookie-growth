package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	R *http.Request
	W http.ResponseWriter
	param map[string]string
}
func (c *Context)ReadJSON(req interface{})error {
	body,err := io.ReadAll(c.R.Body)
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
	data,err:=json.Marshal(resp)
	if err != nil{
		return err
	}
	_,err =c.W.Write(data)
	return err
}
func NewContext(W http.ResponseWriter,R *http.Request,)*Context{
	return &Context {
		R: R,
		W: W,
	}
}
func (c *Context)BadJSON(resp interface{})error{
	return c.WriteJSON(http.StatusInternalServerError,resp)
}