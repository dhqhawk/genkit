package sqlmore

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// sql中支持的数据类型有限,但是它提供了自定义类型的接口
// 如：sql.NullString支持的两种接口
type JsonSql[T any] struct {
	Column T
	Valid  bool // Valid is true if Column is not NULL
}

// Scan 用于将数据库的数据类型转化为go的类型
func (js *JsonSql[T]) Scan(value any) error {
	if value == nil {
		js.Column, js.Valid = *new(T), false
		return nil
	}
	js.Valid = true
	return convertAssign(&js.Column, value)
}

func convertAssign(dest, src any) error {
	return convertAssignRows(dest, src, nil)
}

func convertAssignRows(dest, src any, rows *sql.Rows) error {
	/*	if reflect.ValueOf(dest).IsZero() {
		return errors.New("没有初始化的Column")
	}*/
	switch s := src.(type) {
	case string:
		json.Unmarshal([]byte(s), dest)
	case []byte:
		json.Unmarshal(s, dest)
	}
	return nil
}

// Value 将go的数据转化到数据库中
func (js JsonSql[T]) Value() (driver.Value, error) {
	if !js.Valid {
		return nil, nil
	}
	return json.Marshal(js.Column)

}
