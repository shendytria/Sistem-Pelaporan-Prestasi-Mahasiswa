package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	Repo *repository.LecturerRepository
}

func NewLecturerService(repo *repository.LecturerRepository) *LecturerService {
	return &LecturerService{Repo: repo}
}

func (s *LecturerService) FindByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	return s.Repo.FindByUserID(ctx, userID)
}

func (s *LecturerService) FindAdvisees(ctx context.Context, lecturerID string) ([]map[string]interface{}, error) {
	return s.Repo.FindAdvisees(ctx, lecturerID)
}

func (s *LecturerService) ListHTTP(c *fiber.Ctx) error {
    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Admin" {
        data, err := s.Repo.FindAll(ctx)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(data)
    }

    if role == "Mahasiswa" {
        student, err := s.Repo.FindStudentByUserID(ctx, userID)
        if err != nil || student == nil {
            return c.Status(404).JSON(fiber.Map{"error": "student profile not found"})
        }

        if student.AdvisorID == "" {
            return c.Status(404).JSON(fiber.Map{"error": "advisor not assigned"})
        }

        lecturer, err := s.Repo.FindByID(ctx, student.AdvisorID)
        if err != nil || lecturer == nil {
            return c.Status(404).JSON(fiber.Map{"error": "advisor not found"})
        }

        return c.JSON(lecturer)
    }

    return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
}

func (s *LecturerService) AdviseesHTTP(c *fiber.Ctx) error {
    ctx := context.Background()
    lecturerID := c.Params("id")
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Dosen Wali" {
        myLecturer, err := s.FindByUserID(ctx, userID)
        if err != nil || myLecturer == nil {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
        if myLecturer.ID != lecturerID {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    advisees, err := s.FindAdvisees(ctx, lecturerID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(advisees)
}
