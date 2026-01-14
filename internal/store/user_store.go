package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStore interface {
	Create(*User) error
	GetById(string) (*User, error)
}

func (pg *PostgresUserStore) Create(user *User) error {
	q := `
	INSERT INTO users (email, role, password_hash, is_active)
	VALUES ($1, $2, $3, $4)
	RETURNING id, email, created_at, updated_at
	`
	err := pg.db.
		QueryRow(q, user.Email, user.Role, user.PasswordHash, user.IsActive).
		Scan(
			&user.ID,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresUserStore) GetById(id string) (*User, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	q := `
	SELECT id, email, role, password_hash, is_active, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	usr := &User{}
	err = pg.db.QueryRow(q, id).Scan(
		&usr.ID,
		&usr.Email,
		&usr.Role,
		&usr.PasswordHash,
		&usr.IsActive,
		&usr.CreatedAt,
		&usr.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return usr, nil
}
