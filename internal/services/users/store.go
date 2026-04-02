package users

import (
	"tiny-goclean/internal/types"

	"github.com/jmoiron/sqlx"
)

type userRow struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	IsActive     bool   `db:"is_active"`
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Search() ([]types.User, error) {
	return nil, nil
}

func (s *Store) Create(payload types.CreateUserPayload) error {
	return nil
}

func (s *Store) GetByUsername(username string) types.User {

	q := `
	SELECT * from user
	`
	return nil
}
