package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/jmoiron/sqlx"
)

type completer struct {
	db               *sqlx.DB
	dbType           string
	dbName           string
	keywordsSuggests []prompt.Suggest
	tablesSuggests   []prompt.Suggest
	columnsSuggests  []prompt.Suggest
}

func newCompleter(db *sqlx.DB, dbType, dbName string) *completer {
	return &completer{
		db:     db,
		dbName: dbName,
	}
}

func (c *completer) init() {
	for k := range keywords {
		c.keywordsSuggests = append(c.keywordsSuggests, prompt.Suggest{
			Text:        k,
			Description: "keyword",
		})
	}

	if c.dbType == "mysql" {
		c.initMysql()
	} else {
		c.initPostgres()
	}
}

func (c *completer) initMysql() {
	var tableNames []string
	_ = c.db.Select(&tableNames, `show tables`)

	for _, t := range tableNames {
		c.tablesSuggests = append(c.tablesSuggests, prompt.Suggest{
			Text:        t,
			Description: "table",
		})
	}

	var columns []struct {
		Column string `db:"column_name"`
		Table  string `db:"table_name"`
	}
	_ = c.db.Select(&columns, `SELECT column_name as column_name, table_name as table_name
  FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ?`, c.dbName)

	for _, co := range columns {
		c.columnsSuggests = append(c.columnsSuggests, prompt.Suggest{
			Text:        co.Column,
			Description: co.Table,
		})
	}
}

func (c *completer) initPostgres() {
	var tableNames []string
	_ = c.db.Select(&tableNames, `SELECT tablename
FROM pg_catalog.pg_tables
WHERE schemaname != 'pg_catalog'
AND schemaname != 'information_schema'`)

	for _, t := range tableNames {
		c.tablesSuggests = append(c.tablesSuggests, prompt.Suggest{
			Text:        t,
			Description: "table",
		})
	}

	var columns []struct {
		Column string `db:"column_name"`
		Table  string `db:"table_name"`
	}
	_ = c.db.Select(&columns,
		`SELECT column_name, table_name
FROM information_schema.columns
WHERE table_schema != 'pg_catalog'
AND table_schema != 'information_schema'`)

	for _, co := range columns {
		c.columnsSuggests = append(c.columnsSuggests, prompt.Suggest{
			Text:        co.Column,
			Description: co.Table,
		})
	}
}

func (c *completer) suggest(d prompt.Document) []prompt.Suggest {
	currentLine := d.CurrentLine()
	if currentLine == "" {
		return nil
	}
	fields := strings.Fields(currentLine)
	var lastKeyword string
	for i := len(fields) - 1; i >= 0; i-- {
		if _, ok := keywords[strings.ToUpper(fields[i])]; ok {
			lastKeyword = strings.ToUpper(fields[i])
			break
		}
	}

	var s []prompt.Suggest
	switch lastKeyword {
	case "SELECT", "WHERE":
		s = append(s, c.columnsSuggests...)
	case "FROM":
		s = append(s, c.tablesSuggests...)
	default:
	}

	s = append(s, c.keywordsSuggests...)

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var tableNames []string
var columns []string
