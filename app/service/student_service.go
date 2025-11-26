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
	ctx := context.Background()

	roleIfc := c.Locals("role")
	userIDIfc := c.Locals("user_id")
	if roleIfc == nil || userIDIfc == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	role := roleIfc.(string)
	userID := userIDIfc.(string)

	allStudents, err := s.FindAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var filtered []model.Student

	if role == "Admin" {
		filtered = allStudents
	} else if role == "Mahasiswa" {
		st, _ := s.FindByUserID(ctx, userID)
		if st != nil {
			filtered = append(filtered, *st)
		}
	} else if role == "Dosen Wali" {
		for _, st := range allStudents {
			ok, _ := s.IsMyStudent(ctx, userID, st.ID)
			if ok {
				filtered = append(filtered, st)
			}
		}
	} else {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	return c.JSON(filtered)
}

func (s *StudentService) DetailHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	roleIfc := c.Locals("role")
	userIDIfc := c.Locals("user_id")
	if roleIfc == nil || userIDIfc == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	role := roleIfc.(string)
	userID := userIDIfc.(string)

	st, err := s.FindByID(ctx, id)
	if err != nil || st == nil {
		return c.Status(404).JSON(fiber.Map{"error": "student not found"})
	}

	if role == "Admin" {
		return c.JSON(st)
	} else if role == "Dosen Wali" {
		ok, _ := s.IsMyStudent(ctx, userID, st.ID)
		if ok {
			return c.JSON(st)
		} else {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	} else {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
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
