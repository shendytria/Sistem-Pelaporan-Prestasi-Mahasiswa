package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	Repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
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

func (s *StudentService) ListHTTP(c *fiber.Ctx) error {
	res, err := s.FindAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (s *StudentService) DetailHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := s.FindByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(res)
}

func (s *StudentService) MeHTTP(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	res, err := s.FindByUserID(context.Background(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "student profile not found"})
	}
	return c.JSON(res)
}

func (s *StudentService) UpdateAdvisorHTTP(c *fiber.Ctx) error {
	studentID := c.Params("id")

	type Req struct {
		AdvisorID string `json:"advisorId"`
	}
	var req Req
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
