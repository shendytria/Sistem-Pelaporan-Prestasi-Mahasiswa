package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/middleware"
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

func (s *AchievementMongoService) FindByID(ctx context.Context, id string) (*model.Achievement, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *AchievementMongoService) Update(ctx context.Context, id string, data model.AchievementMongoUpdate) error {
	return s.Repo.Update(ctx, id, data)
}

func (s *AchievementMongoService) SoftDelete(ctx context.Context, id string) error {
	return s.Repo.SoftDelete(ctx, id)
}

func (s *AchievementMongoService) PushAttachment(ctx context.Context, id string, file model.AchievementFile) error {
	return s.Repo.PushAttachment(ctx, id, file)
}

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

func (s *AchievementReferenceService) FindByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	return s.Repo.FindByStudent(ctx, studentID)
}

func (s *AchievementReferenceService) UpdateStatus(ctx context.Context, id string, status string, submittedAt *time.Time, verifiedAt *time.Time, verifiedBy *string, rejectionNote *string) error {
	return s.Repo.UpdateStatus(ctx, id, status, submittedAt, verifiedAt, verifiedBy, rejectionNote)
}

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
		StudentSvc: NewStudentService(studentRepo, refRepo, mongoRepo),
	}
}

func (s *AchievementUsecaseService) List(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "read_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	refs, err := s.RefSvc.FindAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var filtered []model.AchievementReference
	for _, r := range refs {
		switch role {
		case "Admin":
			filtered = append(filtered, r)
		case "Mahasiswa":
			st, _ := s.StudentSvc.FindByUserID(ctx, userID)
			if st != nil && st.ID == r.StudentID {
				filtered = append(filtered, r)
			}
		case "Dosen Wali":
			ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, r.StudentID)
			if ok {
				filtered = append(filtered, r)
			}
		}
	}

	total := len(filtered)
	if offset >= total {
		return c.JSON(fiber.Map{"data": []model.Achievement{}, "page": page, "limit": limit, "total": total, "pages": (total + limit - 1) / limit})
	}

	end := offset + limit
	if end > total {
		end = total
	}
	pageRefs := filtered[offset:end]

	var mongoIDs []string
	for _, r := range pageRefs {
		mongoIDs = append(mongoIDs, r.MongoAchievementID)
	}

	achs, _ := s.MongoSvc.FindMany(ctx, mongoIDs)
	return c.JSON(fiber.Map{"data": achs, "page": page, "limit": limit, "total": total, "pages": (total + limit - 1) / limit})
}

func (s *AchievementUsecaseService) Detail(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "read_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, id)
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

	ach, _ := s.MongoSvc.FindByID(ctx, ref.MongoAchievementID)
	return c.JSON(ach)
}

func (s *AchievementUsecaseService) Create(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "create_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	var studentID string
	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil {
			return c.Status(400).JSON(fiber.Map{"error": "student profile not found"})
		}
		studentID = st.ID
	} else {
		type CreateReq struct {
			StudentID string `json:"studentId"`
		}
		var body CreateReq
		c.BodyParser(&body)
		studentID = body.StudentID
		sid, ok := body.StudentID, body.StudentID != ""
		if !ok || sid == "" {
			return c.Status(400).JSON(fiber.Map{"error": "studentId required"})
		}
		studentID = sid
	}

	var req model.Achievement
	c.BodyParser(&req)
	req.StudentID = studentID

	mongoID, _ := s.MongoSvc.Insert(ctx, &req)
	ref := &model.AchievementReference{
		ID: uuid.New().String(), StudentID: studentID, MongoAchievementID: mongoID.Hex(),
		Status: "draft", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.RefSvc.Insert(ctx, ref)

	return c.JSON(fiber.Map{"message": "achievement created"})
}

func (s *AchievementUsecaseService) Update(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "update_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, err := s.RefSvc.FindByID(ctx, id)
	if err != nil || ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	var req model.AchievementUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	update := model.AchievementMongoUpdate{
		Title:           req.Title,
		Description:     req.Description,
		AchievementType: req.AchievementType,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		UpdatedAt:       time.Now(),
	}

	if update.Title == nil && update.Description == nil && update.AchievementType == nil &&
		update.Details == nil && update.Tags == nil && update.Points == nil {
		return c.Status(400).JSON(fiber.Map{"error": "no fields to update"})
	}

	if err := s.MongoSvc.Update(ctx, ref.MongoAchievementID, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	s.RefSvc.UpdateStatus(ctx, id, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote)

	return c.JSON(fiber.Map{"message": "achievement updated"})
}

func (s *AchievementUsecaseService) Delete(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "delete_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, err := s.RefSvc.FindByID(ctx, id)
	if err != nil || ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	s.MongoSvc.SoftDelete(ctx, ref.MongoAchievementID)
	s.RefSvc.UpdateStatus(ctx, id, "deleted", nil, nil, nil, nil)
	return c.JSON(fiber.Map{"message": "deleted"})
}

func (s *AchievementUsecaseService) Submit(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "update_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, _ := s.RefSvc.FindByID(ctx, id)
	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	now := time.Now()
	s.RefSvc.UpdateStatus(ctx, id, "submitted", &now, nil, nil, nil)
	return c.JSON(fiber.Map{"message": "submitted"})
}

func (s *AchievementUsecaseService) Verify(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "verify_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, _ := s.RefSvc.FindByID(ctx, id)
	if ref.Status == "deleted" {
		return c.Status(400).JSON(fiber.Map{"error": "achievement deleted"})
	}

	if role == "Dosen Wali" && ref.Status != "submitted" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden – final decision already made"})
	}

	if role == "Dosen Wali" {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden – student not your advisee"})
		}
	}

	now := time.Now()
	verifier := userID
	s.RefSvc.UpdateStatus(ctx, id, "verified", nil, &now, &verifier, nil)
	return c.JSON(fiber.Map{"message": "verified"})
}

func (s *AchievementUsecaseService) Reject(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "verify_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	type Req struct {
		Reason string `json:"reason"`
	}
	var req Req
	c.BodyParser(&req)

	ref, _ := s.RefSvc.FindByID(ctx, id)
	if ref.Status == "deleted" {
		return c.Status(400).JSON(fiber.Map{"error": "achievement deleted"})
	}

	if role == "Dosen Wali" && ref.Status != "submitted" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden – final decision already made"})
	}

	if role == "Dosen Wali" {
		ok, _ := s.StudentSvc.IsMyStudent(ctx, userID, ref.StudentID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden – student not your advisee"})
		}
	}

	s.RefSvc.UpdateStatus(ctx, id, "rejected", nil, nil, nil, &req.Reason)
	return c.JSON(fiber.Map{"message": "rejected"})
}

func (s *AchievementUsecaseService) History(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "read_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, _ := s.RefSvc.FindByID(ctx, id)

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

	history := []model.AchievementHistory{
		{Status: "draft", At: &ref.CreatedAt},
	}

	if ref.SubmittedAt != nil {
		history = append(history, model.AchievementHistory{Status: "submitted", At: ref.SubmittedAt})
	}
	if ref.VerifiedAt != nil {
		history = append(history, model.AchievementHistory{Status: "verified", At: ref.VerifiedAt, By: ref.VerifiedBy})
	}
	if ref.Status == "rejected" && ref.RejectionNote != nil {
		history = append(history, model.AchievementHistory{Status: "rejected", At: &ref.UpdatedAt, Note: ref.RejectionNote})
	}
	if ref.Status == "deleted" {
		history = append(history, model.AchievementHistory{Status: "deleted", At: &ref.UpdatedAt})
	}

	return c.JSON(history)
}

func (s *AchievementUsecaseService) AddAttachment(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "update_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.RefSvc.FindByID(ctx, id)
	if err != nil || ref == nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if role == "Mahasiswa" {
		if ref.Status != "draft" && ref.Status != "submitted" {
			return c.Status(400).JSON(fiber.Map{"error": "achievement is finalized"})
		}

		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == "Admin" && ref.Status == "deleted" {
		return c.Status(400).JSON(fiber.Map{"error": "achievement deleted"})
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

	attachment := model.AchievementFile{
		FileName:   file.Filename,
		FileURL:    savePath,
		FileType:   fileType,
		UploadedAt: time.Now(),
	}

	mongoID := ref.MongoAchievementID
	if err := s.MongoSvc.PushAttachment(ctx, mongoID, attachment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update attachment"})
	}

	s.RefSvc.UpdateStatus(ctx, ref.ID, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote)

	return c.JSON(fiber.Map{"message": "attachment added"})
}
