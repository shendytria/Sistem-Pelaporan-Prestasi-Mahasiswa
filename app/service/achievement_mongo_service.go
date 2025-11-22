package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongoService struct {
	Repo *repository.AchievementMongoRepository
}

func NewAchievementMongoService(repo *repository.AchievementMongoRepository) *AchievementMongoService {
	return &AchievementMongoService{Repo: repo}
}

func (s *AchievementMongoService) Insert(ctx context.Context, a *model.Achievement) (primitive.ObjectID, error) {
	return s.Repo.Insert(ctx, a)
}

func (s *AchievementMongoService) FindMany(ctx context.Context, ids []string) ([]model.Achievement, error) {
	return s.Repo.FindMany(ctx, ids)
}
