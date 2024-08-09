package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	dbConnection := os.Getenv("DB_CONNECTION")

	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	migrationsPath := filepath.Join("migrations")

	fmt.Println("Migrations path:", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+migrationsPath,
		"postgres", driver)
	if err != nil {
		log.Fatalln(err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
			return
		}
		log.Fatalln(err)
	}
}
