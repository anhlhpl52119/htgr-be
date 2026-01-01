package store

import (
	"database/sql"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/pressly/goose/v3"

	"fmt"
)

func Open() (*sql.DB, error) {
	fmt.Println("Connected to database...")

	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db: ping %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func MigrateFs(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
