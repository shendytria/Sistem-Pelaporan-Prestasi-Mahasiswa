package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"

	"github.com/gofiber/fiber/v2"
)

type AchievementUsecaseService struct {
	MongoSvc   *AchievementMongoService
	RefSvc     *AchievementReferenceService
	StudentSvc *StudentService
}

func NewAchievementUsecaseService(
	mongoRepo *repository.AchievementMongoRepository,
	refRepo *repository.AchievementReferenceRepository,
	studentRepo *repository.StudentRepository,
) *AchievementUsecaseService {

	return &AchievementUsecaseService{
		MongoSvc:   NewAchievementMongoService(mongoRepo),
		RefSvc:     NewAchievementReferenceService(refRepo),
		StudentSvc: NewStudentService(studentRepo),
	}
}

func (s *AchievementUsecaseService) CreateHTTP(c *fiber.Ctx) error {

	var req model.Achievement

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "studentId required"})
	}

	if err := s.Create(context.Background(), req.StudentID, &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement created"})
}

func (s *AchievementUsecaseService) ListByStudentHTTP(c *fiber.Ctx) error {

	studentID := c.Params("studentId")

	res, err := s.ListByStudent(context.Background(), studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (s *AchievementUsecaseService) Create(ctx context.Context, studentID string, ach *model.Achievement) error {

	mongoID, err := s.MongoSvc.Insert(ctx, ach)
	if err != nil {
		return err
	}

	ref := &model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	return s.RefSvc.Insert(ctx, ref)
}

func (s *AchievementUsecaseService) ListByStudent(ctx context.Context, studentID string) ([]model.Achievement, error) {

	ids, err := s.RefSvc.FindMongoIDsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return s.MongoSvc.FindMany(ctx, ids)
}
