package repository

import (
	"tiny-goclean/internal/model"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func (r *userRepository) FindAll() ([]model.User, error) {
	return nil, nil
}
