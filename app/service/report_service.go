package service

import (
	"context"
    "fmt"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/middleware"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	Repo         repository.ReportRepo
	StudentRepo  repository.StudentRepo
	LecturerRepo repository.LecturerRepo
}

func NewReportService(reportRepo repository.ReportRepo, studentRepo repository.StudentRepo, lecturerRepo  repository.LecturerRepo) *ReportService {
	return &ReportService{
		Repo:         reportRepo,
		StudentRepo:  studentRepo,
		LecturerRepo: lecturerRepo,
	}
}

// Get Achievement Statistics godoc
// @Summary Menampilkan statistik prestasi
// @Description Admin melihat seluruh statistik, Mahasiswa hanya statistik miliknya sendiri, Dosen Wali hanya statistik mahasiswa bimbingannya
// @Security BearerAuth
// @Tags Reports
// @Produce json
// @Success 200 {object} model.StatisticResponse
// @Failure 403 {object} map[string]string
// @Router /reports/statistics [get]
func (s *ReportService) Statistics(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    var filter string

    switch role {
    case "Admin":
        filter = ""

    case "Dosen Wali":
        lecturer, err := s.LecturerRepo.FindByUserID(ctx, userID)
        if err != nil || lecturer == nil {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
        filter = fmt.Sprintf(
            " WHERE student_id IN (SELECT id FROM students WHERE advisor_id = '%s')",
            lecturer.ID,
        )

    case "Mahasiswa":
        student, err := s.StudentRepo.FindByUserID(ctx, userID)
        if err != nil || student == nil {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
        filter = fmt.Sprintf(" WHERE student_id = '%s'", student.ID)

    default:
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    stats, err := s.Repo.GetStatistics(ctx, filter)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(stats)
}

// Get Student Achievement Report godoc
// @Summary Menampilkan statistik prestasi berdasarkan mahasiswa tertentu
// @Description Admin dapat melihat semua mahasiswa, Mahasiswa hanya dirinya sendiri, Dosen Wali hanya mahasiswa bimbingannya
// @Security BearerAuth
// @Tags Reports
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.StudentReportResponse
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reports/student/{id} [get]
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)
    targetStudent := c.Params("id")

    switch role {
    case "Mahasiswa":
        student, err := s.StudentRepo.FindByUserID(ctx, userID)
        if err != nil || student == nil || student.ID != targetStudent {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }

    case "Dosen Wali":
        lecturer, err := s.LecturerRepo.FindByUserID(ctx, userID)
        if err != nil || lecturer == nil {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
        ok, err := s.StudentRepo.CheckAdvisor(ctx, lecturer.ID, targetStudent)
        if err != nil || !ok {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }

    case "Admin":

    default:
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    report, err := s.Repo.GetStudentReport(ctx, targetStudent)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(report)
}
