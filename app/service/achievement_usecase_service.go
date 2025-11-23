package service

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/gofiber/fiber/v2"

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
    ctx := context.Background()
    studentID := c.Params("studentId")

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