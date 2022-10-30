package querylog

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"log"
	"orm_test1"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{}
	db, err := orm_test1.Open("sqlite3", "file:test.db?cache=shared&mode=memory", orm_test1.DBWithMiddleware(
		builder.LogFunc(func(sql string, args []any) {
			fmt.Println(sql)
		}).Build(),
	))
	if err != nil {
		log.Fatal(err)
	}
	_, err = orm_test1.NewSelect[TestModel](db).Get(context.Background())
	assert.Equal(t, nil, err)
}

type TestModel struct {
}
