package main


import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)
var NotFound error = errors.New("data not found")


func Dao(sql_string string)error {
	err := err_demo()
	if err == sql.ErrNoRows {
		return errors.Wrap(NotFound,fmt.Sprintf("data not found sql:%s ",sql_string))
	}
	if err != nil {
		return errors.Wrap(err,fmt.Sprintf("db system err sql:%s",sql_string))
	}

}
func err_demo()error {
	return sql.ErrNoRows
}
func lo()error{
	err := Dao("")
	if errors.Is(err,NotFound) {
		return nil
	}
	if err != nil {

	}
	return
}