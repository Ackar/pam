package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TODO: history
// TODO: config file
// TODO: options
func main() {
	db, err := sqlx.Connect("mysql", "root:pw@tcp(localhost:3306)/classicmodels")
	if err != nil {
		fmt.Printf("Unable to connect to database: %s\n", err)
		os.Exit(1)
	}

	exec := newExecutor(db, &renderer{})
	c := newCompleter(db)
	go c.init()

	p := prompt.New(exec.execute, c.suggest, prompt.OptionWriter(&CustomWriter{}))

	p.Run()
}
