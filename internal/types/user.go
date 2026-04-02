package types

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Search() ([]User, error)
	Create(p CreateUserPayload) error
}

type User struct {
	ID           string         `json:"id"`
	Username     string         `json:"username"`
	PasswordHash hashedPassword `json:"_"`
	IsActive     bool           `json:"is_active"`
}

type hashedPassword string

func (p *hashedPassword) Set(hashed string) error {
	*p = hashedPassword(hashed)
	return nil
}

func (p *hashedPassword) CompareWith(plainText string) (bool, error) {
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

func (p *hashedPassword) GenerateFrom(plainText string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}

	*p = hashedPassword(hashed)
	return nil
}

type CreateUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   *int   `json:"role_id,omitempty"`
}
