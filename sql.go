package main

import (
	"database/sql"
)

var keywords = map[string]interface{}{
	"SELECT":   nil,
	"UPDATE":   nil,
	"INSERT":   nil,
	"FROM":     nil,
	"INTO":     nil,
	"VALUES":   nil,
	"WHERE":    nil,
	"LIMIT":    nil,
	"LIKE":     nil,
	"AND":      nil,
	"OR":       nil,
	"DESCRIBE": nil,
	"JOIN":     nil,
}

func sqlTypeToGo(t *sql.ColumnType) interface{} {
	switch t.DatabaseTypeName() {
	case "VARCHAR", "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT", "UUID",
		"JSON", "JSONB", "DATE", "TIMESTAMP", "TIMESTAMPTZ", "DATETIME", "TIME", "YEAR", "CHAR":
		var s string
		return &s
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT",
		"INT2", "INT4":
		var i int
		return &i
	case "DECIMAL", "FLOAT", "DOUBLE":
		var f float64
		return &f
	case "BOOL":
		var b bool
		return &b
	default:
		var i interface{}
		return &i
	}
}
