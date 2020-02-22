package main

import (
	"database/sql"
	"fmt"
	"os"
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
	case "VARCHAR", "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT":
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
	default:
		fmt.Fprintf(os.Stderr, "unknow type %s\n", t.DatabaseTypeName())
		return nil
	}
}
