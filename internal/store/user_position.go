package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
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
	Update(*UserPosition) error
	GetById(string) (*UserPosition, error)
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

func (pg *PostgresUserPosition) GetById(id string) (*UserPosition, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	q := `
	SELECT id, title, created_at, updated_at
	FROM positions
	WHERE id = $1
	`
	pos := &UserPosition{}
	err = pg.db.QueryRow(q, id).Scan(
		&pos.ID,
		&pos.Title,
		&pos.CreatedAt,
		&pos.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return pos, nil
}
