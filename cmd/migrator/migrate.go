package main

import (
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var  migrationsPath, migrationsTable, pgConnectionString string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&pgConnectionString, "pg-connection-string", "", "PostgreSQL connection string")
	flag.Parse()


	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	if pgConnectionString == "" {
		panic("pg-connection-string is required")
	}

	m, err := migrate.New(
		"file:"+migrationsPath,
		pgConnectionString,
	)
	
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied")
}

type Log struct {
	verbose bool
}

func (l *Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Log) Verbose() bool {
	return false
}