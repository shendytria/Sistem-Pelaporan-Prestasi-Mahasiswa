package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongoRepository struct{}

func NewAchievementMongoRepository() *AchievementMongoRepository {
	return &AchievementMongoRepository{}
}

func (r *AchievementMongoRepository) Insert(ctx context.Context, a *model.Achievement) (primitive.ObjectID, error) {
	a.ID = primitive.NewObjectID()

	_, err := database.Mongo.Collection("achievements").InsertOne(ctx, a)
	return a.ID, err
}

func (r *AchievementMongoRepository) FindMany(ctx context.Context, ids []string) ([]model.Achievement, error) {

	var objIDs []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	cursor, err := database.Mongo.Collection("achievements").
		Find(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Achievement
	for cursor.Next(ctx) {
		var a model.Achievement
		if err := cursor.Decode(&a); err == nil {
			list = append(list, a)
		}
	}

	return list, nil
}
