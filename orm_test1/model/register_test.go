package model

import (
	"github.com/stretchr/testify/assert"
	"orm_test1/internal/errs"
	"reflect"
	"testing"
)

func TestRegister_ParseModel(t *testing.T) {
	testcases := []struct {
		name    string
		input   any
		wantVal *Model
		wantErr error
	}{
		{
			name:    "nil",
			input:   nil,
			wantErr: errs.ErrPointerOnly,
		}, {
			name: "struct",
			input: TestModel{
				Id: 1,
			},
			wantErr: errs.ErrPointerOnly,
		}, {
			name:  "point",
			input: &TestModel{},
			wantVal: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"id": &Field{
						Colname: "id",
						GoName:  "Id",
						Type:    reflect.TypeOf(1),
					},
					"age": &Field{
						Colname: "age",
						GoName:  "Age",
						Type:    reflect.TypeOf(int8(0)),
					},
					"first_name": &Field{
						Colname: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
					},
				},
				Columns: []*Field{&Field{
					Colname: "id",
					GoName:  "Id",
					Type:    reflect.TypeOf(1),
				}, &Field{
					Colname: "age",
					GoName:  "Age",
					Type:    reflect.TypeOf(int8(0)),
				}, &Field{
					Colname: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				}},
			},
		},
	}

	for _, tc := range testcases {
		r := Register{}
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.ParseModel(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, m)
		})
	}
}

type TestModel struct {
	Id        int
	Age       int8
	FirstName string
}
