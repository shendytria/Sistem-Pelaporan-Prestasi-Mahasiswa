package service

import (
	"context"
	"time"
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

func (s *AchievementReferenceService) FindAll(ctx context.Context) ([]model.AchievementReference, error) {
    return s.Repo.FindAll(ctx)
}

func (s *AchievementReferenceService) FindByID(ctx context.Context, id string) (*model.AchievementReference, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *AchievementReferenceService) UpdateStatus(ctx context.Context, id string, status string, verifiedAt *time.Time, verifiedBy *string, rejectionNote *string,) error {
	return s.Repo.UpdateStatus(ctx, id, status, verifiedAt, verifiedBy, rejectionNote)
}
