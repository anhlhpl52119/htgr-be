package users

import (
	"database/sql"
	"errors"
	"tiny-goclean/internal/types"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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

func (s *Store) Create(u types.User, passwordHash string, roleID int) (*types.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var created types.User
	err = tx.QueryRowx(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, is_active",
		u.Username, passwordHash,
	).Scan(&created.ID, &created.Username, &created.IsActive)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
		created.ID, roleID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &created, nil
}

func (s *Store) AssignRole(userID string, roleID int) error {
	_, err := s.db.Exec(
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
		userID, roleID,
	)
	return err
}

func (s *Store) GetByUsername(username string) (*types.User, error) {
	var row userRow
	err := s.db.Get(&row, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &types.User{
		ID:       row.ID,
		Username: row.Username,
		IsActive: row.IsActive,
	}, nil
}

func (s *Store) VerifyPassword(username, plainText string) (*types.User, error) {
	var row userRow
	err := s.db.Get(&row, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(plainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil
		}
		return nil, err
	}

	return &types.User{ID: row.ID, Username: row.Username, IsActive: row.IsActive}, nil
}
