package orm

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"orm/internal/errs"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db := MemoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// From 都不调用
			name: "no from",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM
			name: "with from",
			q:    NewSelector[TestModel](db).From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "empty from",
			q:    NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name: "with db",
			q:    NewSelector[TestModel](db).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewSelector[TestModel](db).From("`test_model_t`").
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			q:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},
		{
			// 非法列
			name:    "invalid column",
			q:       NewSelector[TestModel](db).Where(Not(C("Invalid").GT(18))),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_Get(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	require.NoError(t, err)
	testcases := []struct {
		name     string
		query    string
		mockRows *sqlmock.Rows
		mockErr  error
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:    "single rows",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("deng"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "deng"},
			},
		},
	}

	for _, tc := range testcases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}
	}
	db, err := OpenDB(mockdb)
	require.NoError(t, err)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, res)
		})
	}

}

// func TestSelector_Get(t *testing.T) {
//
//		mockDB, mock, err := sqlmock.New()
//		require.NoError(t, err)
//
//		testCases := []struct {
//			name     string
//			query    string
//			mockErr  error
//			mockRows *sqlmock.Rows
//			wantErr  error
//			wantVal  *TestModel
//		}{
//			{
//				name:    "single row",
//				query:   "SELECT .*",
//				mockErr: nil,
//				mockRows: func() *sqlmock.Rows {
//					rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
//					rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"))
//					return rows
//				}(),
//				wantVal: &TestModel{
//					Id:        123,
//					FirstName: "Ming",
//					Age:       18,
//					LastName:  &sql.NullString{Valid: true, String: "Deng"},
//				},
//			},
//
//			{
//				// SELECT 出来的行数小于你结构体的行数
//				name:    "less columns",
//				query:   "SELECT .*",
//				mockErr: nil,
//				mockRows: func() *sqlmock.Rows {
//					rows := sqlmock.NewRows([]string{"id", "first_name"})
//					rows.AddRow([]byte("123"), []byte("Ming"))
//					return rows
//				}(),
//				wantVal: &TestModel{
//					Id:        123,
//					FirstName: "Ming",
//				},
//			},
//
//			{
//				name:    "invalid columns",
//				query:   "SELECT .*",
//				mockErr: nil,
//				mockRows: func() *sqlmock.Rows {
//					rows := sqlmock.NewRows([]string{"id", "first_name", "gender"})
//					rows.AddRow([]byte("123"), []byte("Ming"), []byte("male"))
//					return rows
//				}(),
//				wantErr: errs.NewErrUnknownColumn("gender"),
//			},
//
//			{
//				name:    "more columns",
//				query:   "SELECT .*",
//				mockErr: nil,
//				mockRows: func() *sqlmock.Rows {
//					rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "first_name"})
//					rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"), []byte("明明"))
//					return rows
//				}(),
//				wantErr: errs.ErrTooManyReturnedColumns,
//			},
//		}
//
//		for _, tc := range testCases {
//			if tc.mockErr != nil {
//				mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
//			} else {
//				mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
//			}
//		}
//
//		db, err := OpenDB(mockDB)
//		require.NoError(t, err)
//		for _, tt := range testCases {
//			t.Run(tt.name, func(t *testing.T) {
//				res, err := NewSelector[TestModel](db).Get(context.Background())
//				assert.Equal(t, tt.wantErr, err)
//				if err != nil {
//					return
//				}
//				assert.Equal(t, tt.wantVal, res)
//			})
//		}
//	}
func MemoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
