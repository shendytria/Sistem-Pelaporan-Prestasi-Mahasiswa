package repository

import (
	"context"
	"time"

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

	now := time.Now()
    a.CreatedAt = now
    a.UpdatedAt = now

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

func (r *AchievementMongoRepository) FindByID(ctx context.Context, id string) (*model.Achievement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var ach model.Achievement
	err = database.Mongo.Collection("achievements").
		FindOne(ctx, bson.M{"_id": oid}).Decode(&ach)

	if err != nil {
		return nil, err
	}

	return &ach, nil
}

func (r *AchievementMongoRepository) Update(ctx context.Context, id string, data map[string]interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	data["updatedAt"] = time.Now()

	_, err = database.Mongo.Collection("achievements").
		UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})

	return err
}

func (r *AchievementMongoRepository) SoftDelete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = database.Mongo.Collection("achievements").
		UpdateOne(ctx, bson.M{"_id": oid},
			bson.M{"$set": bson.M{"deletedAt": time.Now()}},
		)

	return err
}

type AchievementReferenceRepository struct{}

func NewAchievementReferenceRepository() *AchievementReferenceRepository {
	return &AchievementReferenceRepository{}
}

func (r *AchievementReferenceRepository) Insert(ctx context.Context, ref *model.AchievementReference) error {

	const q = `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := database.PG.Exec(ctx, q,
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status,
	)

	return err
}

func (r *AchievementReferenceRepository) FindMongoIDsByStudent(ctx context.Context, studentID string) ([]string, error) {

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

	var list []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			list = append(list, id)
		}
	}

	return list, nil
}

func (r *AchievementReferenceRepository) FindAll(ctx context.Context) ([]model.AchievementReference, error) {

	const q = `
        SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at
        FROM achievement_references
    `

	rows, err := database.PG.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err == nil {
			list = append(list, ref)
		}
	}

	return list, nil
}

func (r *AchievementReferenceRepository) FindByID(ctx context.Context, id string) (*model.AchievementReference, error) {

	const q = `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
			   verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
    `

	row := database.PG.QueryRow(ctx, q, id)

	var ref model.AchievementReference

	err := row.Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func (r *AchievementReferenceRepository) FindByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {

	const q = `
        SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at
        FROM achievement_references
        WHERE student_id::text = $1
        ORDER BY created_at DESC
    `

	rows, err := database.PG.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err == nil {
			list = append(list, ref)
		}
	}

	return list, nil
}

func (r *AchievementReferenceRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
	submittedAt *time.Time,
	verifiedAt *time.Time,
	verifiedBy *string,
	rejectionNote *string,
) error {

	const q = `
        UPDATE achievement_references
        SET status = $1,
            submitted_at = $2,
            verified_at = $3,
            verified_by = $4,
            rejection_note = $5,
            updated_at = NOW()
        WHERE id = $6
    `

	_, err := database.PG.Exec(ctx, q,
		status,
		submittedAt,
		verifiedAt,
		verifiedBy,
		rejectionNote,
		id,
	)

	return err
}
