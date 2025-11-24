package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.Repo.FindByUsername(ctx, username)
}

func (s *UserService) ListHTTP(c *fiber.Ctx) error {
	users, err := s.Repo.FindAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (s *UserService) DetailHTTP(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.Repo.FindByID(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

func (s *UserService) CreateHTTP(c *fiber.Ctx) error {
	var req model.CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	user := model.User{
        ID:           uuid.NewString(),
        Username:     req.Username,
        Email:        req.Email,
        FullName:     req.FullName,
        RoleID:       req.RoleID,
        PasswordHash: req.Password,
		IsActive:     true,
    }

    err := s.Repo.Create(context.Background(), &user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(user)
}

func (s *UserService) UpdateHTTP(c *fiber.Ctx) error {
	id := c.Params("id")

	var req model.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	req.ID = id

	err := s.Repo.Update(context.Background(), &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Updated"})
}

func (s *UserService) DeleteHTTP(c *fiber.Ctx) error {
	id := c.Params("id")

	err := s.Repo.Delete(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted"})
}

func (s *UserService) UpdateRoleHTTP(c *fiber.Ctx) error {
	id := c.Params("id")

	type roleReq struct {
		RoleID string `json:"role_id"`
	}

	var req roleReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	err := s.Repo.UpdateRole(context.Background(), id, req.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role updated"})
}