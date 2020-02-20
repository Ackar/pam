package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

type executor struct {
	db       *sqlx.DB
	renderer *renderer
}

func newExecutor(db *sqlx.DB, renderer *renderer) *executor {
	return &executor{
		db:       db,
		renderer: renderer,
	}
}

func (e *executor) execute(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	if in == "exit" {
		os.Exit(0)
	}
	rows, err := e.db.Query(in)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	var resultRows [][]interface{}
	for rows.Next() {
		arr := buildScanArray(types)
		rows.Scan(arr...)
		resultRows = append(resultRows, arr)
	}

	e.renderer.renderResults(columns, resultRows)
}

func buildScanArray(types []*sql.ColumnType) []interface{} {
	res := make([]interface{}, len(types))
	for i, t := range types {
		res[i] = sqlTypeToGo(t)
	}

	return res
}
