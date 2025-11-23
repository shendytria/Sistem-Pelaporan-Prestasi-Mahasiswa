package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
)

type StudentService struct {
	Repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
}

func (s *StudentService) FindByUserID(ctx context.Context, userID string) (*model.Student, error) {
	return s.Repo.FindByUserID(ctx, userID)
}

func (s *StudentService) IsMyStudent(ctx context.Context, dosenUserID, studentID string) (bool, error) {
    return s.Repo.IsMyStudent(ctx, dosenUserID, studentID)
}