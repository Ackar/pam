package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/jmoiron/sqlx"
)

type executor struct {
	db       *sqlx.DB
	renderer *renderer
	history  *history
}

func newExecutor(db *sqlx.DB, renderer *renderer, history *history) *executor {
	return &executor{
		db:       db,
		renderer: renderer,
		history:  history,
	}
}

func (e *executor) execute(in string) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			ctxCancel()
		}
	}()

	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	if in == "exit" {
		os.Exit(0)
	}

	if in[len(in)-1] != ';' {
		fmt.Println("missing trailing ';'")
		return
	}

	e.history.add(in)

	rows, err := e.db.QueryContext(ctx, in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
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
