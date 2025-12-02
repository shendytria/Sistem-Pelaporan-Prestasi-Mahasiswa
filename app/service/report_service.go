package service

import (
	"context"
    "fmt"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/middleware"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	Repo         *repository.ReportRepository
	StudentRepo  *repository.StudentRepository
	LecturerRepo *repository.LecturerRepository
}

func NewReportService(reportRepo *repository.ReportRepository, studentRepo *repository.StudentRepository, lecturerRepo *repository.LecturerRepository) *ReportService {
	return &ReportService{
		Repo:         reportRepo,
		StudentRepo:  studentRepo,
		LecturerRepo: lecturerRepo,
	}
}

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
