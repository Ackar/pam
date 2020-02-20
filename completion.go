package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/jmoiron/sqlx"
)

type completer struct {
	db               *sqlx.DB
	keywordsSuggests []prompt.Suggest
	tablesSuggests   []prompt.Suggest
	columnsSuggests  []prompt.Suggest
}

func newCompleter(db *sqlx.DB) *completer {
	return &completer{
		db: db,
	}
}

func (c *completer) init() {
	for k := range keywords {
		c.keywordsSuggests = append(c.keywordsSuggests, prompt.Suggest{
			Text:        k,
			Description: "keyword",
		})
	}

	var tableNames []string
	_ = c.db.Select(&tableNames, `show tables`)

	for _, t := range tableNames {
		c.tablesSuggests = append(c.tablesSuggests, prompt.Suggest{
			Text:        t,
			Description: "table",
		})
	}

	var columns []struct {
		Column string `db:"COLUMN_NAME"`
		Table  string `db:"TABLE_NAME"`
	}
	_ = c.db.Select(&columns, `SELECT column_name, table_name
  FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'classicmodels'`)

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
