package homework

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

func InsertStmt(entity interface{}) (string, []interface{}, error) {

	// val := reflect.ValueOf(entity)
	// typ := val.Type()
	// 检测 entity 是否符合我们的要求
	// 我们只支持有限的几种输入

	// 使用 strings.Builder 来拼接 字符串
	// bd := strings.Builder{}

	// 构造 INSERT INTO XXX，XXX 是你的表名，这里我们直接用结构体名字

	// 遍历所有的字段，构造出来的是 INSERT INTO XXX(col1, col2, col3)
	// 在这个遍历的过程中，你就可以把参数构造出来
	// 如果你打算支持组合，那么这里你要深入解析每一个组合的结构体
	// 并且层层深入进去

	// 拼接 VALUES，达成 INSERT INTO XXX(col1, col2, col3) VALUES

	// 再一次遍历所有的字段，要拼接成 INSERT INTO XXX(col1, col2, col3) VALUES(?,?,?)
	// 注意，在第一次遍历的时候我们就已经拿到了参数的值，所以这里就是简单拼接 ?,?,?

	// return bd.String(), args, nil

	var args []interface{}
	if entity == nil {

		return "", args, errInvalidEntity
	}

	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Struct {
		if !(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
			return "", args, errInvalidEntity
		}

	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	fieldNum := typ.NumField()
	if fieldNum == 0 {
		return "", args, errInvalidEntity
	}
	tableName := typ.Name()
	fieldNames := make([]string, 0, fieldNum)
	fields := make([]interface{}, 0, fieldNum)
	var dfs func(p reflect.StructField, value reflect.Value)
	hashmap := make(map[string]bool)
	for i := 0; i < fieldNum; i++ {
		fieldname := typ.Field(i).Name
		if (typ.Field(i).Type.Kind() == reflect.Struct && !CheckInterface(typ.Field(i).Type)) || (typ.Field(i).Type.Kind() == reflect.Ptr && typ.Field(i).Type.Elem().Kind() == reflect.Struct && !CheckInterface(typ.Field(i).Type.Elem())) {
			dfs = func(typ reflect.StructField, val reflect.Value) {
				if typ.Type.Kind() != reflect.Struct || (typ.Type.Kind() == reflect.Struct && CheckInterface(typ.Type)) {
					if typ.Type.Kind() == reflect.Ptr || !(typ.Type.Kind() == reflect.Ptr && typ.Type.Elem().Kind() == reflect.Struct) || (typ.Type.Kind() == reflect.Ptr && typ.Type.Elem().Kind() == reflect.Struct && CheckInterface(typ.Type.Elem())) {
						var field interface{}
						if val.IsZero() {
							field = reflect.Zero(typ.Type).Interface()
						} else {
							field = val.Interface()
						}
						fields = append(fields, field)
						fieldNames = append(fieldNames, typ.Name)
						return
					}
				}

				if _, ok := hashmap[typ.Name]; ok {
					return
				}
				hashmap[typ.Name] = true

				for i := 0; i < typ.Type.NumField(); i++ {
					ftype := typ.Type.Field(i)
					dfs(ftype, val.Field(i))
				}

			}
			dfs(typ.Field(i), val.Field(i))
			continue
		}
		if _, ok := hashmap[typ.Field(i).Name]; ok {
			continue
		}
		hashmap[typ.Field(i).Name] = true
		var field interface{}
		if val.Field(i).IsZero() {
			field = reflect.Zero(typ.Field(i).Type).Interface()

		} else {
			field = val.Field(i).Interface()
		}
		fields = append(fields, field)
		fieldNames = append(fieldNames, fieldname)
	}
	sql := GetInsertSql(tableName, fieldNames)
	return sql, fields, nil

}
func GetInsertSql(tableName string, fieldNames []string) string {
	sqlbuilder := &strings.Builder{}
	sqlbuilder.WriteString(fmt.Sprintf("INSERT INTO `%s`", tableName))
	for index, fieldName := range fieldNames {
		fieldNames[index] = "`" + fieldName + "`"
	}
	sqlbuilder.WriteString(fmt.Sprintf("(%s)", strings.Join(fieldNames, ",")))
	xx := []string{}
	for i := 0; i < len(fieldNames); i++ {
		xx = append(xx, "?")
	}
	sqlbuilder.WriteString(fmt.Sprintf(" VALUES(%s);", strings.Join(xx, ",")))

	return sqlbuilder.String()
}

func CheckInterface(typ reflect.Type) bool {
	if typ.Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem()) {
		return true
	}
	return false
}
