package users

import (
	"database/sql"
	"errors"
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

func (s *Store) GetByUsername(username string) (*types.User, error) {
	var row userRow
	err := s.db.Get(&row, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	u := &types.User{
		ID:       row.ID,
		Username: row.Username,
		IsActive: row.IsActive,
	}
	u.PasswordHash.Set(row.PasswordHash)
	return u, nil
}
