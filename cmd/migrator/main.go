package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

func main() {
	var dbURL, migrationsPath, migrationsTable string

	flag.StringVar(&dbURL, "db-url", "", "PostgreSQL database URL (e.g. postgres://user:password@localhost:5432/dbname?sslmode=disable)")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "schema_migrations", "Name of migrations table")

	flag.Parse()

	if dbURL == "" {
		panic("db-url is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(migrationsPath, fmt.Sprintf("%s&x-migrations-table=%s", dbURL, migrationsTable))

	if err != nil {
		panic("failed to create migrator instance: " + err.Error())
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			os.Exit(0)
			return
		}

		panic(err)
	}

	fmt.Println("Successfully applied migrations")
}
