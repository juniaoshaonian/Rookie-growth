package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"orm/internal/errs"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := MemoryDB(t)
	testCases := []struct {
		name    string
		i       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "single value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			}),
			want: &Query{
				SQL:  "INSERT INTO `test_model` (`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?);",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
		{
			name:    "no value",
			i:       NewInserter[TestModel](db),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			name: "specify columns",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			}).Columns("Age", "FirstName"),
			want: &Query{
				SQL:  "INSERT INTO `test_model` (`age`,`first_name`) VALUES (?,?);",
				Args: []any{int8(18), "Tom"},
			},
		},
		{
			name: "upsert muti",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
				&TestModel{
					Id:        13,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{String: "jerry", Valid: true},
				},
			).OnDuplicateKey().Update(Assign("Age", 19)),
			want: &Query{
				SQL:  "INSERT INTO `test_model` (`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?),(?,?,?,?) ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{String: "Ming", Valid: true}, int64(13), "Tom", int8(18), &sql.NullString{String: "jerry", Valid: true}, 19},
			},
		},
		{
			//
			name: "upsert",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			}).OnDuplicateKey().Update(Assign("Age", 19)),
			want: &Query{
				SQL:  "INSERT INTO `test_model` (`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{String: "Ming", Valid: true}, 19},
			},
		},
		{
			//
			name: "upsert values",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			}).OnDuplicateKey().Update(C("Age")),
			want: &Query{
				SQL:  "INSERT INTO `test_model` (`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `age`=VALUES(`age`);",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.i.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, q, tc.want)
		})
	}
}
