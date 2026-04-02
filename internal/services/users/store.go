package users

import (
	"tiny-goclean/internal/types"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Search() ([]types.User, error) {
	return nil, nil
}
