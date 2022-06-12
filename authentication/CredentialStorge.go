package main
type CredentialStorge interface {
	getPasswordByAppId(appid string)string
}
type CsBasedonMysql struct {
	constr string
	sql string

}