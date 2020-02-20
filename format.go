package main

import (
	"fmt"
	"os"
	"reflect"
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
	if len(rows) == 0 {
		fmt.Println("No rows")
		return
	}

	maxRowSize := 0
	var convertedRows [][]interface{}
	for _, r := range rows {
		convertedRow := convertRow(r)
		convertedRows = append(convertedRows, convertedRow)

		size := rowSize(convertedRow)
		if size > maxRowSize {
			maxRowSize = size
		}
	}

	fmt.Printf("row size is %d\n", maxRowSize)
	if maxRowSize > 100 {
		for _, r := range convertedRows {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleRounded)

			for i, c := range r {
				t.AppendRow(table.Row{columns[i], c})
			}

			t.Render()
		}

		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	var header []interface{}
	for _, c := range columns {
		header = append(header, c)
	}
	t.AppendHeader(header)
	for _, r := range convertedRows {
		t.AppendRow(r)
	}
	t.Render()
}

func rowSize(r []interface{}) int {
	var buf strings.Builder
	for _, c := range r {
		buf.WriteString(fmt.Sprint(c))
	}

	return buf.Len()
}

func convertRow(r []interface{}) []interface{} {
	var row []interface{}
	for _, a := range r {
		row = append(row, ptrToType(a))
	}

	return row
}

func ptrToType(v interface{}) interface{} {
	value := reflect.ValueOf(v)
	if value.IsNil() {
		return nil
	}
	return value.Elem().Interface()
}
