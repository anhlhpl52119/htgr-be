package service

import "tiny-goclean/internal/model"

type UserRepository interface {
	FindAll() ([]model.User, error)
}

type UserService struct {
	repo UserRepository
}
