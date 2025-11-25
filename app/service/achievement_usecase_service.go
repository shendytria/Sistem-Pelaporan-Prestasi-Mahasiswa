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

	now := time.Now()
	err = s.RefSvc.UpdateStatus(ctx, achID, "submitted", &now, nil, nil, nil)
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

	err = s.RefSvc.UpdateStatus(ctx, achID, "verified", nil, &now, &verifier, nil)
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

	err = s.RefSvc.UpdateStatus(ctx, achID, "rejected", nil, nil, nil, &rejection)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "rejected"})
}

func (s *AchievementUsecaseService) DetailHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleMahasiswa {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == constant.RoleDosenWali {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	ach, err := s.MongoSvc.FindByID(ctx, ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mongo not found"})
	}

	return c.JSON(ach)
}

func (s *AchievementUsecaseService) UpdateHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	userID := c.Locals("user_id").(string)
	role := c.Locals("role_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleMahasiswa {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	err = s.MongoSvc.Update(ctx, ref.MongoAchievementID, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "updated"})
}

func (s *AchievementUsecaseService) DeleteHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	if role == constant.RoleMahasiswa {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	err = s.MongoSvc.SoftDelete(ctx, ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	err = s.RefSvc.UpdateStatus(ctx, achID, "deleted", nil, nil, nil, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "pg update failed"})
	}

	return c.JSON(fiber.Map{"message": "deleted"})
}

func (s *AchievementUsecaseService) HistoryHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleMahasiswa {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == constant.RoleDosenWali {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	type HistoryItem struct {
		Status string     `json:"status"`
		At     time.Time  `json:"at"`
		By     *string    `json:"by,omitempty"`
		Note   *string    `json:"note,omitempty"`
	}

	history := []HistoryItem{
		{
			Status: "draft",
			At:     ref.CreatedAt,
		},
	}

	if ref.SubmittedAt != nil {
		history = append(history, HistoryItem{
			Status: "submitted",
			At:     *ref.SubmittedAt,
		})
	}

	if ref.VerifiedAt != nil {
		history = append(history, HistoryItem{
			Status: "verified",
			At:     *ref.VerifiedAt,
			By:     ref.VerifiedBy,
		})
	}

	if ref.Status == "rejected" && ref.RejectionNote != nil {
		history = append(history, HistoryItem{
			Status: "rejected",
			At:     ref.UpdatedAt,
			Note:   ref.RejectionNote,
		})
	}

	if ref.Status == "deleted" {
		history = append(history, HistoryItem{
			Status: "deleted",
			At:     ref.UpdatedAt,
		})
	}

	return c.JSON(history)
}

func (s *AchievementUsecaseService) AddAttachmentHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == constant.RoleMahasiswa {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	type AttachReq struct {
		FileName string `json:"fileName"`
		FileURL  string `json:"fileUrl"`
		FileType string `json:"fileType"`
	}

	var req AttachReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if req.FileName == "" || req.FileURL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "fileName and fileUrl are required"})
	}

	ach, err := s.MongoSvc.FindByID(ctx, ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mongo not found"})
	}

	newFile := model.AchievementFile{
		FileName:   req.FileName,
		FileURL:    req.FileURL,
		FileType:   req.FileType,
		UploadedAt: time.Now(),
	}

	ach.Attachments = append(ach.Attachments, newFile)

	err = s.MongoSvc.Update(ctx, ref.MongoAchievementID, map[string]interface{}{
		"attachments": ach.Attachments,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "attachment added"})
}