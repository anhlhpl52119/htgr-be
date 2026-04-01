package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type DB struct {
	Client *sqlx.DB
}

const maxOpenConn = 10
const maxIdleConn = 5
const maxLifeTime = 5 * time.Minute

func NewConnection(dns string) (*DB, error) {
	conn := &DB{}

	db, err := sqlx.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	err = checkConnection(db)
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxIdleTime(maxIdleConn)
	db.SetConnMaxLifetime(maxLifeTime)
	db.SetMaxOpenConns(maxOpenConn)
	conn.Client = db

	return conn, nil
}

func checkConnection(db *sqlx.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("db ping: %w", err)
	}
	fmt.Println("*** Ping db successfully **")
	return nil
}

func CheckMigration(db *sqlx.DB, migrationDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	vers, err := goose.GetDBVersion(db.DB)
	if err != nil {
		return fmt.Errorf("get migration version: %w", err)
	}

}
