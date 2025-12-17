package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/middleware"
)

type AchievementService struct {
	Repo       repository.AchievementRepo
	StudentSvc StudentServiceInterface
}

func NewAchievementService(repo repository.AchievementRepo, studentSvc StudentServiceInterface) *AchievementService {
	return &AchievementService{Repo: repo, StudentSvc: studentSvc}
}

// List Achievements godoc
// @Summary Menampilkan daftar prestasi
// @Description Admin melihat semua, Mahasiswa hanya miliknya, Dosen Wali hanya milik mahasiswa bimbingannya
// @Security BearerAuth
// @Tags Achievements
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements [get]
func (s *AchievementService) List(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "read_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	refs, err := s.Repo.FindAllReferences(ctx)
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

	achs, _ := s.Repo.FindManyMongo(ctx, mongoIDs)
	return c.JSON(fiber.Map{"data": achs, "page": page, "limit": limit, "total": total, "pages": (total + limit - 1) / limit})
}

// Detail Achievement godoc
// @Summary Melihat detail 1 prestasi
// @Description Admin melihat semua, Mahasiswa hanya miliknya, Dosen Wali hanya milik mahasiswa bimbingannya
// @Security BearerAuth
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} model.Achievement
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id} [get]
func (s *AchievementService) Detail(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "read_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.Repo.FindReferenceByID(ctx, id)
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

	ach, _ := s.Repo.FindByIDMongo(ctx, ref.MongoAchievementID)
	return c.JSON(ach)
}

// Create Achievement godoc
// @Summary Membuat prestasi baru
// @Description Admin dapat membuat untuk student lain, Mahasiswa membuat untuk dirinya sendiri
// @Security BearerAuth
// @Tags Achievements
// @Accept json
// @Produce json
// @Param request body model.Achievement true "Achievement body"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements [post]
func (s *AchievementService) Create(c *fiber.Ctx) error {
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
		if body.StudentID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "studentId required"})
		}
		studentID = body.StudentID
	}

	var req model.Achievement
	c.BodyParser(&req)
	req.StudentID = studentID

	mongoID, _ := s.Repo.InsertMongo(ctx, &req)
	ref := &model.AchievementReference{
		ID: uuid.New().String(), StudentID: studentID, MongoAchievementID: mongoID,
		Status: "draft", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.Repo.InsertReference(ctx, ref)

	return c.JSON(fiber.Map{"message": "achievement created"})
}

// Update Achievement godoc
// @Summary Update prestasi (hanya status DRAFT)
// @Description Admin bisa override student_id, Mahasiswa dapat mengedit miliknya selama masih berstatus draft
// @Security BearerAuth
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param request body model.AchievementUpdate true "Partial update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [put]
func (s *AchievementService) Update(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "update_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.Repo.FindReferenceByID(ctx, id)
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

	if role == "Admin" && req.StudentID != nil && *req.StudentID != ref.StudentID {
		ref.StudentID = *req.StudentID
		now := time.Now()
		s.Repo.UpdateStatus(ctx, ref.ID, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote, req.StudentID)
		_ = s.Repo.UpdateMongo(ctx, ref.MongoAchievementID, &model.AchievementMongoUpdate{
			StudentID: req.StudentID,
			UpdatedAt: now,
		})
		_ = s.Repo.InsertReference(ctx, &model.AchievementReference{
			ID:                 ref.ID,
			StudentID:          ref.StudentID,
			MongoAchievementID: ref.MongoAchievementID,
			Status:             ref.Status,
			CreatedAt:          ref.CreatedAt,
			UpdatedAt:          now,
		})
	}

	old, _ := s.Repo.FindByIDMongo(ctx, ref.MongoAchievementID)
	if old == nil {
		old = &model.Achievement{}
	}

	update := model.AchievementMongoUpdate{
		UpdatedAt: time.Now(),
	}

	if req.Title != nil {
		update.Title = req.Title
	} else {
		v := old.Title
		update.Title = &v
	}

	if req.Description != nil {
		update.Description = req.Description
	} else {
		v := old.Description
		update.Description = &v
	}

	if req.AchievementType != nil {
		update.AchievementType = req.AchievementType
	} else {
		v := old.AchievementType
		update.AchievementType = &v
	}

	if req.Details != nil {
		update.Details = req.Details
	} else {
		v := old.Details
		update.Details = &v
	}

	if req.Tags != nil {
		update.Tags = req.Tags
	} else {
		update.Tags = old.Tags
	}

	if req.Points != nil {
		update.Points = req.Points
	} else {
		v := old.Points
		update.Points = &v
	}

	if err := s.Repo.UpdateMongo(ctx, ref.MongoAchievementID, &update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	s.Repo.UpdateStatus(ctx, id, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote, nil)

	return c.JSON(fiber.Map{"message": "achievement updated"})
}

// Delete Achievement godoc
// @Summary Menghapus prestasi (soft delete)
// @Description Admin bisa menghapus mana saja, Mahasiswa hanya boleh menghapus yang berstatus draft miliknya sendiri
// @Security BearerAuth
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [delete]
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "delete_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.Repo.FindReferenceByID(ctx, id)
	if err != nil || ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	if role == "Mahasiswa" {
		st, _ := s.StudentSvc.FindByUserID(ctx, userID)
		if st == nil || st.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	s.Repo.SoftDeleteMongo(ctx, ref.MongoAchievementID)
	s.Repo.UpdateStatus(ctx, id, "deleted", nil, nil, nil, nil, nil)
	return c.JSON(fiber.Map{"message": "deleted"})
}

// Submit Achievement godoc
// @Summary Submit prestasi dari draft → submitted
// @Description Admin bisa submit untuk student lain, Mahasiswa hanya boleh submit miliknya sendiri
// @Security BearerAuth
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/submit [post]
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "update_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, _ := s.Repo.FindReferenceByID(ctx, id)
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
	s.Repo.UpdateStatus(ctx, id, "submitted", &now, nil, nil, nil, nil)
	return c.JSON(fiber.Map{"message": "submitted"})
}

// Verify Achievement godoc
// @Summary Verifikasi prestasi mahasiswa
// @Description Admin bisa verifikasi mana saja, Dosen Wali hanya bisa verifikasi milik mahasiswanya & hanya status submitted
// @Security BearerAuth
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/verify [post]
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "verify_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	ref, _ := s.Repo.FindReferenceByID(ctx, id)
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
	s.Repo.UpdateStatus(ctx, id, "verified", nil, &now, &verifier, nil, nil)
	return c.JSON(fiber.Map{"message": "verified"})
}

// Reject Achievement godoc
// @Summary Menolak prestasi dengan alasan tertentu
// @Description Admin bisa menolak mana saja, Dosen wali hanya menolak milik mahasiswanya & hanya status submitted
// @Security BearerAuth
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param request body map[string]string true "Alasan penolakan"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/reject [post]
func (s *AchievementService) Reject(c *fiber.Ctx) error {
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

	ref, _ := s.Repo.FindReferenceByID(ctx, id)
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

	s.Repo.UpdateStatus(ctx, id, "rejected", nil, nil, nil, &req.Reason, nil)
	return c.JSON(fiber.Map{"message": "rejected"})
}

// Achievement History godoc
// @Summary Mendapatkan riwayat perubahan status prestasi
// @Description Admin bisa melihat riwayat perubahan status prestasi mana saja, Mahasiswa hanya bisa melihat riwayat prestasi miliknya, Dosen wali hanya bisa melihat riwayat prestasi milik mahasiswa bimbingannya
// @Security BearerAuth
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {array} model.AchievementHistory
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/history [get]
func (s *AchievementService) History(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "read_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, _ := s.Repo.FindReferenceByID(ctx, id)

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

	history := []model.AchievementHistory{{Status: "draft", At: &ref.CreatedAt}}
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

// Add Attachment godoc
// @Summary Menambahkan lampiran pada prestasi dalam bentuk JSON (URL file) atau upload file langsung
// @Description Admin bisa menambah file mana saja, Mahasiswa hanya bisa menambah pada prestasi miliknya yang berstatus draft atau submitted
// @Security BearerAuth
// @Tags Achievements
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param file formData file false "Upload file"
// @Param fileName formData string false "Nama file (jika JSON)"
// @Param fileUrl formData string false "URL file (jika JSON)"
// @Param fileType formData string false "Tipe file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id}/attachments [post]
func (s *AchievementService) AddAttachment(c *fiber.Ctx) error {
	if !middleware.HasPermission(c, "update_achievement") {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	ctx := context.Background()
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.Repo.FindReferenceByID(ctx, id)
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

	var jsonInput struct {
		FileName string `json:"fileName"`
		FileURL  string `json:"fileUrl"`
		FileType string `json:"fileType"`
	}
	if err := c.BodyParser(&jsonInput); err == nil && jsonInput.FileName != "" && jsonInput.FileURL != "" {
		if jsonInput.FileType == "" {
			jsonInput.FileType = "unknown"
		}

		attachment := model.AchievementFile{
			FileName:   jsonInput.FileName,
			FileURL:    jsonInput.FileURL,
			FileType:   jsonInput.FileType,
			UploadedAt: time.Now(),
		}
		if err := s.Repo.PushAttachmentMongo(ctx, ref.MongoAchievementID, attachment); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to add attachment to mongo",
			})
		}
		_ = s.Repo.UpdateStatus(ctx, ref.ID, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote, nil)

		return c.JSON(fiber.Map{"message": "attachment added"})
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
	if err := s.Repo.PushAttachmentMongo(ctx, ref.MongoAchievementID, attachment); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to add attachment to mongo",
		})
	}
	_ = s.Repo.UpdateStatus(ctx, ref.ID, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote, nil)

	return c.JSON(fiber.Map{"message": "attachment added"})
}
