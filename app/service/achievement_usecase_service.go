package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
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

func (s *AchievementUsecaseService) ListByStudentHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	studentID := c.Params("id")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	if role == "Dosen Wali" {
		ok, err := s.StudentSvc.IsMyStudent(ctx, userID, studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	ids, err := s.RefSvc.FindMongoIDsByStudent(ctx, studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	achs, err := s.MongoSvc.FindMany(ctx, ids)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(achs)
}

func (s *AchievementUsecaseService) ListHTTP(c *fiber.Ctx) error {
	ctx := context.Background()

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	refs, err := s.RefSvc.FindAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var filteredRefs []model.AchievementReference

	for _, r := range refs {
		if role == "Mahasiswa" {
			st, _ := s.StudentSvc.FindByUserID(ctx, userID)
			if st != nil && st.ID == r.StudentID {
				filteredRefs = append(filteredRefs, r)
			}
		} else if role == "Dosen Wali" {
			ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, r.StudentID)
			if ok {
				filteredRefs = append(filteredRefs, r)
			}
		} else if role == "Admin" {
			filteredRefs = append(filteredRefs, r)
		}
	}

	var ids []string
	for _, r := range filteredRefs {
		ids = append(ids, r.MongoAchievementID)
	}

	achs, err := s.MongoSvc.FindMany(ctx, ids)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(achs)
}

func (s *AchievementUsecaseService) DetailHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == "Dosen Wali" {
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

func (s *AchievementUsecaseService) CreateHTTP(c *fiber.Ctx) error {
	ctx := context.Background()

	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "Mahasiswa" && role != "Admin" {
		return c.Status(403).JSON(fiber.Map{"error": "only mahasiswa or admin can create achievement"})
	}

	var studentID string

	if role == "Mahasiswa" {
		st, err := s.StudentSvc.FindByUserID(ctx, userID)
		if err != nil || st == nil {
			return c.Status(400).JSON(fiber.Map{"error": "student data not found for this user"})
		}
		studentID = st.ID
	} else if role == "Admin" {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
		}
		if val, ok := body["studentId"].(string); ok && val != "" {
			studentID = val
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "studentId is required for admin"})
		}
	}

	var req model.Achievement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	req.StudentID = studentID

	mongoID, err := s.MongoSvc.Insert(ctx, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ref := &model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.RefSvc.Insert(ctx, ref); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement created"})
}

func (s *AchievementUsecaseService) UpdateHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	var req model.Achievement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	ach, err := s.MongoSvc.FindByID(ctx, ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mongo achievement not found"})
	}

	ach.Title = req.Title
	ach.Description = req.Description
	ach.AchievementType = req.AchievementType
	ach.Details = req.Details
	ach.Tags = req.Tags
	ach.Points = req.Points

	updateData := map[string]interface{}{
		"title":           ach.Title,
		"description":     ach.Description,
		"achievementType": ach.AchievementType,
		"details":         ach.Details,
		"tags":            ach.Tags,
		"points":          ach.Points,
	}

	if err := s.MongoSvc.Update(ctx, ref.MongoAchievementID, updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement updated"})
}

func (s *AchievementUsecaseService) DeleteHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	if role == "Mahasiswa" {
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

func (s *AchievementUsecaseService) SubmitHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == "Mahasiswa" {
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

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	if role == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == "Dosen Wali" {
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

	role := c.Locals("role").(string)
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

	if role == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if role == "Dosen Wali" {
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

func (s *AchievementUsecaseService) HistoryHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	achID := c.Params("id")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == "Dosen Wali" {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	type HistoryItem struct {
		Status string    `json:"status"`
		At     time.Time `json:"at"`
		By     *string   `json:"by,omitempty"`
		Note   *string   `json:"note,omitempty"`
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

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, achID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file is required"})
	}

	fileType := c.FormValue("fileType")
	if fileType == "" {
		fileType = "unknown"
	}

	savePath := "./uploads/" + file.Filename
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
	}

	ach, err := s.MongoSvc.FindByID(ctx, ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mongo not found"})
	}

	newFile := model.AchievementFile{
		FileName:   file.Filename,
		FileURL:    savePath,
		FileType:   fileType,
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
