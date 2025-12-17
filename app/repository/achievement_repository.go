package repository

import (
	"context"
	"time"
	"errors"

	"prestasi_mhs/app/model"
	"prestasi_mhs/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRepo interface {
	InsertMongo(ctx context.Context, a *model.Achievement) (string, error)
	FindManyMongo(ctx context.Context, ids []string) ([]model.Achievement, error)
	FindByIDMongo(ctx context.Context, id string) (*model.Achievement, error)
	UpdateMongo(ctx context.Context, id string, data *model.AchievementMongoUpdate) error
	SoftDeleteMongo(ctx context.Context, id string) error
	PushAttachmentMongo(ctx context.Context, id string, file model.AchievementFile) error
	InsertReference(ctx context.Context, ref *model.AchievementReference) error
	FindAllReferences(ctx context.Context) ([]model.AchievementReference, error)
	FindReferenceByID(ctx context.Context, id string) (*model.AchievementReference, error)
	UpdateStatus(ctx context.Context, id, status string, submittedAt, verifiedAt *time.Time, verifiedBy, note *string, studentID *string) error
	FindMongoIDsByStudent(ctx context.Context, studentID string) ([]string, error)
}

type AchievementRepository struct{}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{}
}

func (r *AchievementRepository) InsertMongo(ctx context.Context, a *model.Achievement) (string, error) {
	a.ID = primitive.NewObjectID()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err := database.Mongo.Collection("achievements").InsertOne(ctx, a)
	return a.ID.Hex(), err
}

func (r *AchievementRepository) FindManyMongo(ctx context.Context, ids []string) ([]model.Achievement, error) {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, oid)
		}
	}
	cursor, err := database.Mongo.Collection("achievements").Find(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
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

func (r *AchievementRepository) FindByIDMongo(ctx context.Context, id string) (*model.Achievement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var ach model.Achievement
	err = database.Mongo.Collection("achievements").
		FindOne(ctx, bson.M{"_id": oid}).Decode(&ach)
	return &ach, err
}

func (r *AchievementRepository) UpdateMongo(ctx context.Context, id string, data *model.AchievementMongoUpdate) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = database.Mongo.Collection("achievements").
		UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})
	return err
}

func (r *AchievementRepository) SoftDeleteMongo(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = database.Mongo.Collection("achievements").
		UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"deletedAt": time.Now()}})
	return err
}

func (r *AchievementRepository) PushAttachmentMongo(ctx context.Context, id string, file model.AchievementFile,) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := database.Mongo.Collection("achievements").
		UpdateOne(
			ctx,
			bson.M{"_id": oid},
			bson.M{
				"$push": bson.M{"attachments": file},
				"$set":  bson.M{"updatedAt": time.Now()},
			},
		)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("mongo achievement not found")
	}

	return nil
}

func (r *AchievementRepository) InsertReference(ctx context.Context, ref *model.AchievementReference) error {
	const q = `
		INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`
	_, err := database.PG.Exec(ctx, q, ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status)
	return err
}

func (r *AchievementRepository) FindAllReferences(ctx context.Context) ([]model.AchievementReference, error) {
	const q = `SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at FROM achievement_references`
	rows, err := database.PG.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &ref.CreatedAt, &ref.UpdatedAt)
		list = append(list, ref)
	}
	return list, nil
}

func (r *AchievementRepository) FindReferenceByID(ctx context.Context, id string) (*model.AchievementReference, error) {
	const q = `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by,
		       rejection_note, created_at, updated_at
        FROM achievement_references WHERE id = $1
    `
	row := database.PG.QueryRow(ctx, q, id)

	var ref model.AchievementReference
	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)
	return &ref, err
}

func (r *AchievementRepository) UpdateStatus(ctx context.Context, id, status string, submittedAt, verifiedAt *time.Time, verifiedBy, note *string, studentID *string) error {
	const q = `
		UPDATE achievement_references
		SET status=$1, submitted_at=$2, verified_at=$3, verified_by=$4, rejection_note=$5, updated_at=NOW(), student_id = COALESCE($6, student_id)
		WHERE id=$7
	`
	_, err := database.PG.Exec(ctx, q, status, submittedAt, verifiedAt, verifiedBy, note, studentID, id)
	return err
}

func (r *AchievementRepository) FindMongoIDsByStudent(ctx context.Context, studentID string) ([]string, error) {
	const q = `
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE student_id::text = $1
	`

	rows, err := database.PG.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
