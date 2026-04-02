package types

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Search() ([]User, error)
}

type User struct {
	ID           string         `json:"id" db:"id"`
	Username     string         `json:"username" db:"username"`
	PasswordHash hashedPassword `json:"oo" db:"password_hash"`
	IsActive     bool           `json:"is_active" db:"is_active"`
}

type hashedPassword string

func (p *hashedPassword) Set(plainText string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}

	*p = hashedPassword(hashed)
	return nil
}

func (p *hashedPassword) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(*p), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err //internal server error
		}
	}

	return true, nil
}
