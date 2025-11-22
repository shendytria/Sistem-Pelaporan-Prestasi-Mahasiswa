package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
)

type LecturerService struct {
	Repo *repository.LecturerRepository
}

func NewLecturerService(repo *repository.LecturerRepository) *LecturerService {
	return &LecturerService{Repo: repo}
}

func (s *LecturerService) FindByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	return s.Repo.FindByUserID(ctx, userID)
}
