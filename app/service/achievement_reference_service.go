package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
)

type AchievementReferenceService struct {
	Repo *repository.AchievementReferenceRepository
}

func NewAchievementReferenceService(repo *repository.AchievementReferenceRepository) *AchievementReferenceService {
	return &AchievementReferenceService{Repo: repo}
}

func (s *AchievementReferenceService) Insert(ctx context.Context, ref *model.AchievementReference) error {
	return s.Repo.Insert(ctx, ref)
}

func (s *AchievementReferenceService) FindMongoIDsByStudent(ctx context.Context, studentID string) ([]string, error) {
	return s.Repo.FindMongoIDsByStudent(ctx, studentID)
}
