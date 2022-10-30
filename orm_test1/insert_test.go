package orm_test1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := memoryDB(t)
	testcases := []struct {
		name      string
		i         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			//插入单个值得全部列，其实就是插入一行
			name: "single value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
			}),
			wantQuery: &Query{
				Sql:  "INSERT INTO `test_model` (`id`,`age`,`first_name`)  VALUES (?,?,?);",
				Args: []any{12, int8(18), "Tom"},
			},
		},
		{
			name: "muti value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
			}, &TestModel{
				Id:        11,
				FirstName: "daming",
				Age:       19,
			}),
			wantQuery: &Query{
				Sql:  "INSERT INTO `test_model` (`id`,`age`,`first_name`)  VALUES (?,?,?),(?,?,?);",
				Args: []any{12, int8(18), "Tom", 11, int8(19), "daming"},
			},
		},
		{
			name: "specify column",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
			}).Columns("Id", "Age"),
			wantQuery: &Query{
				Sql:  "INSERT INTO `test_model` (`id`,`age`)  VALUES (?,?);",
				Args: []any{12, int8(18)},
			},
		},
		{
			name: "upsert",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
			}).OnDuplicatekey().Update(Assign("Age", 19)),
			wantQuery: &Query{
				Sql:  "INSERT INTO `test_model` (`id`,`age`,`first_name`)  VALUES (?,?,?) ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{12, int8(18), "Tom", 19},
			},
		},
		{
			name: "upsert",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
			}).OnDuplicatekey().Update(Assign("Age", 19), C("FirstName")),
			wantQuery: &Query{
				Sql:  "INSERT INTO `test_model` (`id`,`age`,`first_name`)  VALUES (?,?,?) ON DUPLICATE KEY UPDATE `age`=?,`first_name`=VALUES(`first_name`);",
				Args: []any{12, int8(18), "Tom", 19},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.i.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}
