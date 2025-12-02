package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
    "prestasi_mhs/middleware"

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

func (s *LecturerService) FindAdvisees(ctx context.Context, lecturerID string) ([]model.Advisee, error) {
	return s.Repo.FindAdvisees(ctx, lecturerID)
}

func (s *LecturerService) List(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Mahasiswa" {
        student, err := s.Repo.FindStudentByUserID(ctx, userID)
        if err != nil || student == nil {
            return c.Status(404).JSON(fiber.Map{"error": "student profile not found"})
        }

        lecturer, err := s.Repo.FindByID(ctx, student.AdvisorID)
        if err != nil || lecturer == nil {
            return c.Status(404).JSON(fiber.Map{"error": "advisor not found"})
        }

        return c.JSON(fiber.Map{
            "data":  []model.Lecturer{*lecturer},
            "page":  1,
            "limit": 1,
            "total": 1,
            "pages": 1,
        })
    }

    data, err := s.Repo.FindAll(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    total := len(data)
    if offset >= total {
        return c.JSON(fiber.Map{
            "data":  []model.Lecturer{},
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

    return c.JSON(fiber.Map{
        "data":  data[offset:end],
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}

func (s *LecturerService) Advisees(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "verify_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    lecturerID := c.Params("id")
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Dosen Wali" {
        myLecturer, err := s.FindByUserID(ctx, userID)
        if err != nil || myLecturer == nil || myLecturer.ID != lecturerID {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    advisees, err := s.FindAdvisees(ctx, lecturerID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    total := len(advisees)
    if offset >= total {
        return c.JSON(fiber.Map{
            "data":  []model.Advisee{},
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

    return c.JSON(fiber.Map{
        "data":  advisees[offset:end],
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}