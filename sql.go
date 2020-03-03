package main

import (
	"database/sql"
	"fmt"
	"os"
)

var keywords = map[string]interface{}{
	"SELECT": nil,
	"UPDATE": nil,
	"INSERT": nil,
	"FROM":   nil,
	"INTO":   nil,
	"VALUES": nil,
	"WHERE":  nil,
	"LIMIT":  nil,
	"LIKE":   nil,
	"AND":    nil,
	"OR":     nil,
}

func sqlTypeToGo(t *sql.ColumnType) interface{} {
	switch t.DatabaseTypeName() {
	case "VARCHAR", "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT", "UUID",
		"JSON", "JSONB", "DATE", "TIMESTAMP", "TIMESTAMPTZ", "DATETIME", "TIME":
		var s string
		return &s
	case "INT", "DECIMAL", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT",
		"INT2", "INT4":
		var i int
		return &i
	case "FLOAT", "DOUBLE":
		var f float64
		return &f
	case "CHAR":
		var r rune
		return &r
	case "BOOL":
		var b bool
		return &b
	default:
		fmt.Fprintf(os.Stderr, "unknow type %s\n", t.DatabaseTypeName())
		return nil
	}
}
