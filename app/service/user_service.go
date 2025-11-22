package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.Repo.FindByUsername(ctx, username)
}
