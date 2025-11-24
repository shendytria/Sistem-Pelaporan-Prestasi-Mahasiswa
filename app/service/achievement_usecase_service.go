package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/constant"
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
	ctx := context.Background()

	userID := c.Locals("user_id").(string)
	role := c.Locals("role_id").(string)

	if role != constant.RoleMahasiswa {
		return c.Status(403).JSON(fiber.Map{"error": "only mahasiswa can create achievement"})
	}

	st, err := s.StudentSvc.FindByUserID(ctx, userID)
	if err != nil || st == nil {
		return c.Status(400).JSON(fiber.Map{"error": "student data not found for this user"})
	}

	var req model.Achievement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	req.StudentID = st.ID

	if err := s.Create(ctx, st.ID, &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement created"})
}

func (s *AchievementUsecaseService) ListByStudentHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	studentID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	if role == constant.RoleMahasiswa {
		myData, err := s.StudentSvc.FindByUserID(ctx, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if myData.ID != studentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == constant.RoleDosenWali {
		isMyStudent, err := s.StudentSvc.IsMyStudent(ctx, userID, studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if !isMyStudent {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	res, err := s.ListByStudent(ctx, studentID)
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

func (s *AchievementUsecaseService) ListAllHTTP(c *fiber.Ctx) error {

	res, err := s.ListAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (s *AchievementUsecaseService) ListAll(ctx context.Context) ([]model.Achievement, error) {

	refs, err := s.RefSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, r := range refs {
		ids = append(ids, r.MongoAchievementID)
	}

	return s.MongoSvc.FindMany(ctx, ids)
}

func (s *AchievementUsecaseService) ListMineHTTP(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Locals("user_id").(string)

	st, err := s.StudentSvc.FindByUserID(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.ListByStudent(ctx, st.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (s *AchievementUsecaseService) SubmitHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleMahasiswa {
		student, err := s.StudentSvc.FindByUserID(ctx, userID)
		if err != nil || student == nil {
			return c.Status(403).JSON(fiber.Map{"error": "not a student"})
		}

		if student.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	err = s.RefSvc.UpdateStatus(ctx, achID, "submitted", nil, nil, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "submitted"})
}

func (s *AchievementUsecaseService) VerifyHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	if role == constant.RoleMahasiswa {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleDosenWali {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "must be submitted first"})
	}

	now := time.Now()
	verifier := userID

	err = s.RefSvc.UpdateStatus(ctx, achID, "verified", &now, &verifier, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "verified"})
}

func (s *AchievementUsecaseService) RejectHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	type Req struct {
		Reason string `json:"reason"`
	}
	var req Req
	c.BodyParser(&req)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleMahasiswa {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if role == constant.RoleDosenWali {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	rejection := req.Reason

	err = s.RefSvc.UpdateStatus(ctx, achID, "rejected", nil, nil, &rejection)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "rejected"})
}
