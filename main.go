package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/c-bata/go-prompt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/ssh/terminal"
)

type dbConfig struct {
	Type string
	DSN  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s DB_NAME\n", os.Args[0])
		os.Exit(1)
	}
	var config map[string]dbConfig

	confFile, err := os.Open(os.Getenv("HOME") + "/.pam.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open config file: %s\n", err)
		os.Exit(1)
	}
	defer confFile.Close()

	err = json.NewDecoder(confFile).Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't decode config file: %s\n", err)
		os.Exit(1)
	}

	dbName := os.Args[1]
	dbRegexp := regexp.MustCompile(dbName)
	var matches []string

	for db, _ := range config {
		if dbRegexp.MatchString(db) {
			matches = append(matches, db)
		}
	}

	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "could not match %q to a configured db\n", dbName)
		os.Exit(1)
	} else if len(matches) > 1 {
		fmt.Fprintf(os.Stderr, "%q matched too many configued dbs\n", dbName)
		for _, db := range matches {
			fmt.Fprintf(os.Stderr, "%s\n", db)
		}
		os.Exit(1)
	}

	matchedDB := matches[0]
	conf, found := config[matchedDB]
	if !found {
		fmt.Fprintf(os.Stderr, "could not find configuration for %q\n", matchedDB)
		os.Exit(1)
	}

	var schemaName string
	if conf.Type == "mysql" {
		cfg, err := mysql.ParseDSN(conf.DSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid DSN: %s\n", err)
			os.Exit(1)
		}
		schemaName = cfg.DBName
	} else if conf.Type == "postgres" {
		cfg, err := pgx.ParseConfig(conf.DSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid DSN: %s\n", err)
			os.Exit(1)
		}
		schemaName = cfg.Database
	}

	var driver string
	switch conf.Type {
	case "mysql":
		driver = "mysql"
	case "postgres":
		driver = "pgx"
	default:
		fmt.Fprintf(os.Stderr, "unknown database type %q\n", conf.Type)
		os.Exit(1)
	}
	db, err := sqlx.Connect(driver, conf.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %s\n", err)
		os.Exit(1)
	}

	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		h := newHistory()
		defer h.save()

		exec := newExecutor(db,
			newRenderer(int(prompt.NewStandardInputParser().GetWinSize().Col)),
			h,
		)
		c := newCompleter(db, conf.Type, schemaName)
		go c.init()

		p := prompt.New(exec.execute, c.suggest,
			prompt.OptionWriter(&CustomWriter{}),
			prompt.OptionPrefix(matchedDB+"> "),
			prompt.OptionHistory(h.load()),
		)

		p.Run()
	} else {
		exec := newExecutor(db,
			newRenderer(80),
			newHistory(),
		)
		query, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read input: %s\n", err)
		}
		exec.execute(string(query) + ";")
	}
}
