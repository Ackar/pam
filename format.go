package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/jedib0t/go-pretty/table"
)

type CustomWriter struct {
	prompt.PosixWriter
}

func (c *CustomWriter) displayWord(w string) {
	switch wordType(w) {
	case keyword:
		c.PosixWriter.SetColor(prompt.Blue, prompt.DefaultColor, false)
	case literal:
		c.PosixWriter.SetColor(prompt.Red, prompt.DefaultColor, false)
	case number:
		c.PosixWriter.SetColor(prompt.DarkGreen, prompt.DefaultColor, false)
	default:
	}
	c.PosixWriter.WriteStr(w)
	c.PosixWriter.SetColor(prompt.DefaultColor, prompt.DefaultColor, false)
}

func (c *CustomWriter) WriteStr(s string) {
	current := ""
	inQuotes := false
	for i, r := range s {
		if r == '\'' && (i == 0 || (i > 0 && s[i-1] != '\\')) {
			if inQuotes {
				current += string(r)
				c.displayWord(current)
				current = ""
				inQuotes = false
				continue
			}
			inQuotes = true
		}

		if !inQuotes && (r == ' ' || r == ';' || r == ',') {
			c.displayWord(current)
			c.PosixWriter.WriteStr(string(r))
			current = ""
			continue
		}

		current += string(r)
	}

	c.displayWord(current)
}

type tokenType int

const (
	keyword tokenType = iota
	literal
	number
	other
)

func wordType(w string) tokenType {
	if len(w) == 0 {
		return other
	}

	if _, ok := keywords[strings.ToUpper(w)]; ok {
		return keyword
	}

	if w[0] == '\'' && w[len(w)-1] == '\'' {
		return literal
	}

	_, err := strconv.ParseInt(w, 10, 64)
	if err == nil {
		return number
	}

	return other
}

type renderer struct {
}

func (r *renderer) renderResults(columns []string, rows [][]interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	var header []interface{}
	for _, c := range columns {
		header = append(header, c)
	}
	t.AppendHeader(header)
	for _, r := range rows {
		var row []interface{}
		for _, a := range r {
			switch v := a.(type) {
			case *string:
				row = append(row, *v)
			case *int:
				row = append(row, *v)
			}
		}
		t.AppendRow(row)
	}
	t.Render()
}
