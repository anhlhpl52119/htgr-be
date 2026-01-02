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
	Phone        string    `json:"phone"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStore interface {
	CreateUser(*User) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	q := `
	INSERT INTO users (email, phone, username, password_hash, bio)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, email, created_at, updated_at
	`
	err := pg.db.
		QueryRow(q, user.Email, user.Phone, user.Username, user.PasswordHash, user.Bio).
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

func (pg *PostgresUserStore) GetUserByID(id string) error {
	return nil
}
