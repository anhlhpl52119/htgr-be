package store

import (
	"database/sql"
	"time"
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
	User(*User) error
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
	return nil, nil
}
