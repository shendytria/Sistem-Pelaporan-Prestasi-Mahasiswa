package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
    "prestasi_mhs/middleware"

	"github.com/gofiber/fiber/v2"
)

type StudentServiceInterface interface {
    FindByUserID(ctx context.Context, userID string) (*model.Student, error)
    IsMyStudent(ctx context.Context, dosenUserID, studentID string) (bool, error)
}


type StudentService struct {
    Repo     repository.StudentRepo
    AchievementRepo repository.AchievementRepo
}

func NewStudentService(
    studentRepo repository.StudentRepo,
    achievementRepo repository.AchievementRepo,
) *StudentService {
    return &StudentService{
        Repo:     studentRepo,
        AchievementRepo: achievementRepo,
    }
}

func (s *StudentService) FindByUserID(ctx context.Context, userID string) (*model.Student, error) {
	return s.Repo.FindByUserID(ctx, userID)
}

func (s *StudentService) IsMyStudent(ctx context.Context, dosenUserID, studentID string) (bool, error) {
    return s.Repo.IsMyStudent(ctx, dosenUserID, studentID)
}

func (s *StudentService) FindAll(ctx context.Context) ([]model.Student, error) {
	return s.Repo.FindAll(ctx)
}

func (s *StudentService) FindByID(ctx context.Context, id string) (*model.Student, error) {
	return s.Repo.FindByID(ctx, id)
}

// List Students godoc
// @Summary Menampilkan daftar mahasiswa
// @Description Admin melihat semua mahasiswa, Dosen wali hanya melihat mahasiswa bimbingannya.
// @Security BearerAuth
// @Tags Students
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /students [get]
func (s *StudentService) List(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    students, err := s.Repo.FindAll(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    var filtered []model.Student
    switch role {
    case "Admin":
        filtered = students
    case "Dosen Wali":
        for _, st := range students {
            ok, _ := s.Repo.IsMyStudent(ctx, userID, st.ID)
            if ok {
                filtered = append(filtered, st)
            }
        }
    default:
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    total := len(filtered)
    if offset >= total {
        return c.JSON(fiber.Map{"data": []model.Student{}, "page": page, "limit": limit, "total": total, "pages": (total + limit - 1) / limit})
    }

    end := offset + limit
    if end > total {
        end = total
    }

    return c.JSON(fiber.Map{
        "data":  filtered[offset:end],
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}

// Detail Student godoc
// @Summary Mendapatkan detail 1 mahasiswa
// @Description Admin dapat melihat semua mahasiswa, Dosen wali hanya dapat melihat mahasiswa bimbingannya, Mahasiswa hanya dapat melihat dirinya sendiri
// @Security BearerAuth
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.Student
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /students/{id} [get]
func (s *StudentService) Detail(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    id := c.Params("id")
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    st, err := s.Repo.FindByID(ctx, id)
    if err != nil || st == nil {
        return c.Status(404).JSON(fiber.Map{"error": "student not found"})
    }

    if role == "Mahasiswa" {
        me, _ := s.Repo.FindByUserID(ctx, userID)
        if me == nil || me.ID != st.ID {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    if role == "Dosen Wali" {
        ok, _ := s.Repo.IsMyStudent(ctx, userID, st.ID)
        if !ok {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    return c.JSON(st)
}

// Student Achievements godoc
// @Summary Menampilkan daftar prestasi milik mahasiswa tertentu
// @Description Admin dapat melihat semua, Mahasiswa hanya bisa melihat prestasinya sendiri, Dosen wali hanya bisa melihat prestasi mahasiswa bimbingannya
// @Security BearerAuth
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students/{id}/achievements [get]
func (s *StudentService) ListByStudent(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    studentID := c.Params("id")
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Mahasiswa" {
        me, _ := s.Repo.FindByUserID(ctx, userID)
        if me == nil || me.ID != studentID {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    if role == "Dosen Wali" {
        ok, _ := s.Repo.IsMyStudent(ctx, userID, studentID)
        if !ok {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    ids, err := s.AchievementRepo.FindMongoIDsByStudent(ctx, studentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    total := len(ids)

    if offset >= total {
        return c.JSON(fiber.Map{
            "data":  []model.Achievement{},
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + limit - 1) / limit,
        })
    }

    end := offset + limit
    if end > total {
        end = total
    }

    achs, err := s.AchievementRepo.FindManyMongo(ctx, ids[offset:end])
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "data":  achs,
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}

// Update Advisor godoc
// @Summary Mengubah dosen pembimbing mahasiswa
// @Description Hanya Admin yang dapat mengubah dosen wali seorang mahasiswa
// @Security BearerAuth
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body map[string]string true "Advisor update body"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students/{id}/advisor [put]
func (s *StudentService) UpdateAdvisor(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "manage_users") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    studentID := c.Params("id")

    var req struct {
        AdvisorID string `json:"advisorId"`
    }
    c.BodyParser(&req)

    if req.AdvisorID == "" {
        return c.Status(400).JSON(fiber.Map{"error": "advisorId required"})
    }

    err := s.Repo.UpdateAdvisor(context.Background(), studentID, req.AdvisorID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "advisor updated"})
}