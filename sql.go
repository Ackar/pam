package main

import (
	"database/sql"
	"fmt"
)

var keywords = map[string]interface{}{
	"SELECT": nil,
	"UPDATE": nil,
	"FROM":   nil,
	"VALUES": nil,
	"WHERE":  nil,
	"LIMIT":  nil,
	"LIKE":   nil,
	"AND":    nil,
	"OR":     nil,
}

func sqlTypeToGo(t *sql.ColumnType) interface{} {
	switch t.DatabaseTypeName() {
	case "VARCHAR":
		var s string
		return &s
	case "INT", "DECIMAL":
		var i int
		return &i
	default:
		fmt.Printf("unknow type %s\n", t.DatabaseTypeName())
		return nil
	}
}
