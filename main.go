package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

// TODO: history
// TODO: options

type dbConfig struct {
	Type string
	DSN  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s DB_NAME\n", os.Args[0])
		os.Exit(2)
	}
	var config map[string]dbConfig

	confFile, err := os.Open(os.Getenv("HOME") + "/.pam.json")
	if err != nil {
		fmt.Printf("couldn't open config file: %s\n", err)
		os.Exit(2)
	}
	defer confFile.Close()

	err = json.NewDecoder(confFile).Decode(&config)
	if err != nil {
		fmt.Printf("couldn't decode config file: %s\n", err)
		os.Exit(3)
	}

	dbName := os.Args[1]
	conf, found := config[dbName]
	if !found {
		fmt.Printf("could not find configuration for %q\n", dbName)
		os.Exit(4)
	}

	db, err := sqlx.Connect(conf.Type, conf.DSN)
	if err != nil {
		fmt.Printf("Unable to connect to database: %s\n", err)
		os.Exit(1)
	}

	exec := newExecutor(db, &renderer{})
	c := newCompleter(db)
	go c.init()

	p := prompt.New(exec.execute, c.suggest,
		prompt.OptionWriter(&CustomWriter{}),
		prompt.OptionPrefix(dbName+"> "),
		prompt.OptionPrefixTextColor(prompt.DarkGray),
	)

	p.Run()
}
