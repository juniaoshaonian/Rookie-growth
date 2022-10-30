package orm_test1

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"orm_test1/internal/errs"
	"orm_test1/internal/valuer"
	"testing"
)

func TestSelect_Where(t *testing.T) {
	db := memoryDB(t)
	testcases := []struct {
		name    string
		Builder QueryBuilder
		wantVal *Query
		wantErr error
	}{
		{
			//查全部
			name:    "any",
			Builder: NewSelect[TestModel](db),
			wantVal: &Query{
				Sql: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 查询单个条件
			name:    "single where",
			Builder: NewSelect[TestModel](db).Where([]*Predicate{C("Age").EQ(int8(12))}),
			wantVal: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE `age`=?;",
				Args: []any{int8(12)},
			},
		}, {
			name:    "and",
			Builder: NewSelect[TestModel](db).Where([]*Predicate{C("Age").EQ(int8(12)).And(C("FirstName").EQ("daming"))}),
			wantVal: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE  (`age`=?) AND (`first_name`=?) ;",
				Args: []any{int8(12), "daming"},
			},
		},
		{
			name:    "or",
			Builder: NewSelect[TestModel](db).Where([]*Predicate{C("Age").EQ(int8(12)).Or(C("FirstName").EQ("daming"))}),
			wantVal: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE  (`age`=?) OR (`first_name`=?) ;",
				Args: []any{int8(12), "daming"},
			},
		},
		{
			// 多个条件
			name:    "MUTI",
			Builder: NewSelect[TestModel](db).Where([]*Predicate{C("Age").EQ(int8(12)), (C("FirstName").EQ("daming"))}),
			wantVal: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE  (`age`=?) AND (`first_name`=?) ;",
				Args: []any{int8(12), "daming"},
			},
		},
		{
			name:    "NOT",
			Builder: NewSelect[TestModel](db).Where([]*Predicate{C("Age").EQ(int8(12)).Not()}),
			wantVal: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE NOT (`age`=?) ;",
				Args: []any{int8(12)},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.Builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, q)
		})
	}
}

func TestSelect_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	testcases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:    "single row",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age"})
				rows.AddRow([]byte("123"), []byte("ming"), []byte("18"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "ming",
				Age:       int8(18),
			},
		},
		{
			name:    "no rows",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age"})
				return rows
			}(),
			wantErr: ErrNoRows,
		},
		//{
		//	name:    "invalid row",
		//	query:   "SELECT .*",
		//	mockErr: nil,
		//	mockRows: func() *sqlmock.Rows {
		//		rows := sqlmock.NewRows([]string{"id", "first_name", "age", "gender"})
		//		rows.AddRow([]byte("123"), []byte("ming"), []byte("18"), []byte("male"))
		//		return rows
		//	}(),
		//	wantVal: &TestModel{
		//		Id:        123,
		//		FirstName: "ming",
		//		Age:       int8(18),
		//	},
		//	wantErr: errs.NewErrUnknownColumn("gender"),
		//},
	}
	for _, tc := range testcases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}

	}
	db, err := OpenDB(mockDB, DBWithCreator(valuer.NewUnsafevaluer))
	require.NoError(t, err)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelect[TestModel](db).Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}

}
func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

func TestSelect_Select(t *testing.T) {
	db := memoryDB(t)
	testcases := []struct {
		name    string
		Builder QueryBuilder
		wantVal *Query
		wantErr error
	}{
		{
			// 查询具体字段
			name:    "simple columns",
			Builder: NewSelect[TestModel](db).Select(C("Age")),
			wantVal: &Query{
				Sql: "SELECT `age` FROM `test_model`;",
			},
		},
		{
			//查询聚合函数
			name:    "Aggregate",
			Builder: NewSelect[TestModel](db).Select(Max("Age")),
			wantVal: &Query{
				Sql: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},

		{
			//聚合函数的别名
			name:    "alias",
			Builder: NewSelect[TestModel](db).Select(Max("Age").Alias("max_age")),
			wantVal: &Query{
				Sql: "SELECT MAX(`age`) AS `max_age` FROM `test_model`;",
			},
		},
		{
			// 普通列的别名
			name:    "alias",
			Builder: NewSelect[TestModel](db).Select(C("Age").Alias("max_age")),
			wantVal: &Query{
				Sql: "SELECT `age` AS `max_age` FROM `test_model`;",
			},
		}, {
			// rowexpression
			name:    "rowexpr",
			Builder: NewSelect[TestModel](db).Select(Raw("COUNT (DISTINT `first_name`)")),
			wantVal: &Query{
				Sql: "SELECT COUNT (DISTINT `first_name`) FROM `test_model`;",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.Builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, q)
		})
	}
}

type TestModel struct {
	Id        int
	Age       int8
	FirstName string
}

func TestSelect_GetMulti(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	require.NoError(t, err)
	testcases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantVal  []*TestModel
		wantErr  error
	}{
		{
			name:    "muti rows",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name"})
				rows.AddRow([]byte("123"), []byte("18"), []byte("xiaoming"))
				rows.AddRow([]byte("124"), []byte("16"), []byte("daming"))
				return rows
			}(),
			wantVal: []*TestModel{
				&TestModel{
					Id:        123,
					FirstName: "xiaoming",
					Age:       int8(18),
				},
				&TestModel{
					Id:        124,
					FirstName: "daming",
					Age:       int8(16),
				},
			},
		},
		{
			name:    "zero row",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name"})
				return rows
			}(),
			wantErr: ErrNoRows,
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
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelect[TestModel](db).GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}

}

func TestSelector_Join(t *testing.T) {
	db := memoryDB(t)

	type Order struct {
		Id        int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId  int

		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 虽然泛型是 Order，但是我们传入 OrderDetail
			name: "specify table",
			q:    NewSelect[Order](db).From(TableOf(&OrderDetail{})),
			wantQuery: &Query{
				Sql: "SELECT * FROM `order_detail`;",
			},
		},
		{
			name: "join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{})
				return NewSelect[Order](db).
					From(t1.Join(t2).On(t1.C("Id").EQ(t2.C("OrderId"))))
			}(),
			wantQuery: &Query{
				Sql: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` ON `t1`.`id` = `order_id`);",
			},
		},
		{
			name: "multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelect[Order](db).
					From(t1.Join(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						Join(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				Sql: "SELECT * FROM ((`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "left multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelect[Order](db).
					From(t1.LeftJoin(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						LeftJoin(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				Sql: "SELECT * FROM ((`order` AS `t1` LEFT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) LEFT JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "right multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")

				return NewSelect[Order](db).
					From(t1.RightJoin(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						RightJoin(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				Sql: "SELECT * FROM ((`order` AS `t1` RIGHT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) RIGHT JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},

		{
			name: "join multiple using",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{})
				return NewSelect[Order](db).
					From(t1.Join(t2).Using("UsingCol1", "UsingCol2"))
			}(),
			wantQuery: &Query{
				Sql: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` USING (`using_col1`,`using_col2`));",
			},
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

func TestSelector_SubQuery(t *testing.T) {
	db := memoryDB(t)
	type Order struct {
		Id        int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId  int

		UsingCol1 string
		UsingCol2 string
	}
	testcases := []struct {
		name    string
		input   QueryBuilder
		wantVal *Query
		wantErr error
	}{
		{
			// 测试from类型的子查询
			name: "from",
			input: func() QueryBuilder {
				sub := NewSelect[OrderDetail](db).AsSub("t1")
				return NewSelect[Order](db).From(sub)
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM (SELECT * FROM `order_detail`) AS `t1`;",
			},
		},
		{
			name: "in",
			input: func() QueryBuilder {
				sub := NewSelect[Order](db).Select(C("Id")).AsSub("sub")
				return NewSelect[OrderDetail](db).Where([]*Predicate{C("OrderId").In(sub)})
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM `order_detail` WHERE `order_id` IN (SELECT `id` FROM `order`);",
			},
		},
		{
			name: "Exist",
			input: func() QueryBuilder {
				sub := NewSelect[Order](db).Select(C("Id")).AsSub("sub")
				return NewSelect[OrderDetail](db).Where([]*Predicate{Exist(sub)})
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM `order_detail` WHERE  EXIST (SELECT `id` FROM `order`);",
			},
		},
		{
			name: "Not EXIST",
			input: func() QueryBuilder {
				sub := NewSelect[Order](db).Select(C("Id")).AsSub("sub")
				return NewSelect[OrderDetail](db).Where([]*Predicate{Exist(sub).Not()})
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM `order_detail` WHERE  NOT  ( EXIST (SELECT `id` FROM `order`)) ;",
			},
		},
		{
			name: "some any all",
			input: func() QueryBuilder {
				sub := NewSelect[Order](db).Select(C("Id")).AsSub("sub")
				return NewSelect[OrderDetail](db).Where([]*Predicate{C("OrderId").EQ(ANY(sub))})
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM `order_detail` WHERE `order_id` = ANY (SELECT `id` FROM `order`);",
			},
		},
		{
			name: "some and any",
			input: func() QueryBuilder {
				sub := NewSelect[Order](db).Select(C("Id")).AsSub("sub")
				return NewSelect[OrderDetail](db).Where([]*Predicate{C("OrderId").GT(SOME(sub)), C("OrderId").LT(ANY(sub))})
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM `order_detail` WHERE  (`order_id` > SOME (SELECT `id` FROM `order`))  AND  (`order_id` < ANY (SELECT `id` FROM `order`)) ;",
			},
		},
		{
			// join查询套子查询
			name: "join and sub",
			input: func() QueryBuilder {
				sub := NewSelect[OrderDetail](db).AsSub("sub")
				return NewSelect[Order](db).From(TableOf(&Order{}).Join(sub).On(C("Id").EQ(sub.C("OrderId")))).Select(sub.C("ItemId"))
			}(),
			wantVal: &Query{
				Sql: "SELECT `sub`.`item_id` FROM (`order` JOIN (SELECT * FROM `order_detail`) AS `sub` ON `id` = `sub`.`order_id`);",
			},
		},
		{
			name: "table and left join",
			input: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelect[OrderDetail](db).AsSub("sub")
				return NewSelect[Order](db).From(sub.Join(t1).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM ((SELECT * FROM `order_detail`) AS `sub` JOIN `order` ON `id` = `sub`.`order_id`);",
			},
		},
		{
			name: "join and join",
			input: func() QueryBuilder {
				sub1 := NewSelect[OrderDetail](db).AsSub("sub1")
				sub2 := NewSelect[OrderDetail](db).AsSub("sub2")
				return NewSelect[Order](db).From(sub1.RightJoin(sub2).Using("Id"))
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM ((SELECT * FROM `order_detail`) AS `sub1` RIGHT JOIN (SELECT * FROM `order_detail`) AS `sub2` USING (`id`));",
			},
		},
		{
			name: "join sub sub",
			input: func() QueryBuilder {
				sub1 := NewSelect[OrderDetail](db).AsSub("sub1")
				sub2 := NewSelect[OrderDetail](db).From(sub1).AsSub("sub2")
				t1 := TableOf(&Order{}).As("o1")
				return NewSelect[Order](db).From(sub2.Join(t1).Using("Id"))
			}(),
			wantVal: &Query{
				Sql: "SELECT * FROM ((SELECT * FROM (SELECT * FROM `order_detail`) AS `sub1`) AS `sub2` JOIN `order` AS `o1` USING (`id`));",
			},
		},
		{
			name: "invalid field",
			input: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelect[OrderDetail](db).AsSub("sub")
				return NewSelect[Order](db).Select(sub.C("Invalid")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantErr: errs.NewErrField("Invalid"),
		},
		{
			name: "invalid field in predicates",
			input: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelect[OrderDetail](db).AsSub("sub")
				return NewSelect[Order](db).Select(sub.C("ItemId")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("Invalid"))))
			}(),
			wantErr: errs.NewErrField("Invalid"),
		},
		{
			name: "invalid field in aggregate function",
			input: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelect[OrderDetail](db).AsSub("sub")
				return NewSelect[Order](db).Select(Max("Invalid")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantErr: errs.NewErrField("Invalid"),
		},
		{
			name: "not selected",
			input: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelect[OrderDetail](db).Select(C("OrderId")).AsSub("sub")
				return NewSelect[Order](db).Select(sub.C("ItemId")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantErr: errs.NewErrField("ItemId"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.input.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, q)
		})
	}
}
