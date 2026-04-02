package users

import (
	"tiny-goclean/internal/types"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store types.UserStore
}

func NewService(store types.UserStore) *Service {
	return &Service{store: store}
}

func (s *Service) Search() ([]types.User, error) {
	return s.store.Search()
}

func (s *Service) CreateUser(payload types.CreateUserPayload) (*types.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		return nil, err
	}

	roleID := types.DefaultRoleID
	if payload.RoleID != nil {
		roleID = *payload.RoleID
	}

	return s.store.Create(
		types.User{Username: payload.Username},
		string(hashed),
		roleID,
	)
}
