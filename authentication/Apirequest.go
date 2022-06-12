package main

import "time"

type Apirequest struct {
	Appid string
	Baseurl string
	token string
	create_time time.Duration
}

func (a *Apirequest)getBaseUrl()string{
	return a.Baseurl
}
func (a *Apirequest)getToken()string{
	return a.token
}
func (a *Apirequest)getAppId()string {
	return a.Appid
}
