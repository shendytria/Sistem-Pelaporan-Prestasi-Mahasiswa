package service

import (
	"context"
	"prestasi_mhs/app/repository"

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

func (s *ReportService) StatisticsHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	var filter string
	var args []interface{}

	if role == "Admin" {
		filter = ""

	} else if role == "Dosen Wali" {
		lecturer, err := s.LecturerRepo.FindByUserID(ctx, userID)
		if err != nil || lecturer == nil {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		filter = " WHERE student_id IN (SELECT id FROM students WHERE advisor_id = $1)"
		args = append(args, lecturer.ID)

	} else if role == "Mahasiswa" {
		student, err := s.StudentRepo.FindByUserID(ctx, userID)
		if err != nil || student == nil {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		filter = " WHERE student_id = $1"
		args = append(args, student.ID)

	} else {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	stats, err := s.Repo.GetStatistics(ctx, filter, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}

func (s *ReportService) StudentReportHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)
	targetStudent := c.Params("id")

	if role == "Mahasiswa" {
		student, err := s.StudentRepo.FindByUserID(ctx, userID)
		if err != nil || student == nil || student.ID != targetStudent {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

	} else if role == "Dosen Wali" {
		lecturer, err := s.LecturerRepo.FindByUserID(ctx, userID)
		if err != nil || lecturer == nil {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

		ok, err := s.StudentRepo.CheckAdvisor(ctx, lecturer.ID, targetStudent)
		if err != nil || !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

	} else if role == "Admin" {

	} else {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	report, err := s.Repo.GetStudentReport(ctx, targetStudent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(report)
}
