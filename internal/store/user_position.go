package store

import (
	"database/sql"
	"time"
)

type PostgresUserPosition struct {
	db *sql.DB
}

type UserPosition struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" `
}

func NewUserPosition(pgDb *sql.DB) *PostgresUserPosition {
	return &PostgresUserPosition{
		db: pgDb,
	}
}

type UserPositionStore interface {
	Create(*UserPosition) error
}

func (pg *PostgresUserPosition) Create(pos *UserPosition) error {
	q := `
	INSERT INTO positions (title)
	VALUES $1
	RETURNING id, created_at, updated_at
	`
	err := pg.db.QueryRow(q, pos.Title).Scan(
		&pos.ID,
		&pos.CreatedAt,
		&pos.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
