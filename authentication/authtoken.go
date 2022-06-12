package main

import "time"

type AuthToken struct {
	token string
	createTime time.Duration
	expiredTimeInterval time.Duration
}

func (a *AuthToken)isExpired()bool {
	var consuming_time time.Duration
	consuming_time = time.Duration(time.Now().Unix()) - a.createTime
	if consuming_time > a.expiredTimeInterval {
		return false
	}
	return true
}
func (a *AuthToken)match(authtoken AuthToken)bool{
	if authtoken.token == a.token {
		return true
	}
	return false
}
func (a *AuthToken)GetToken()string{
	return a.token
}