package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostgresPosition struct {
	db *sql.DB
}

type Position struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" `
}

func NewPosition(pgDb *sql.DB) *PostgresPosition {
	return &PostgresPosition{
		db: pgDb,
	}
}

type PositionStore interface {
	Create(*Position) error
	Update(Position) error
	GetById(string) (*Position, error)
}

func (pg *PostgresPosition) Create(pos *Position) error {
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

func (pg *PostgresPosition) GetById(id string) (*Position, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	q := `
	SELECT id, title, created_at, updated_at
	FROM positions
	WHERE id = $1
	`
	pos := &Position{}
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

func (pg *PostgresPosition) Update(pos Position) error {
	_, err := uuid.Parse(pos.ID)
	if err != nil {
		return errors.New("Invalid id format")
	}

	q := `
	UPDATE positions
	SET title = $1,
	WHERE id = $2
	`
	_, err = pg.db.Exec(q, pos.Title, pos.ID)

	if err != nil {
		return err
	}

	return nil
}
