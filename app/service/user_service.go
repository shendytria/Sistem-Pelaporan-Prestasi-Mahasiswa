package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/utils"

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

func (s *UserService) List(c *fiber.Ctx) error {
	ctx := context.Background()

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	users, total, err := s.Repo.FindAll(ctx, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  users,
		"page":  page,
		"limit": limit,
		"total": total,
		"pages": (total + limit - 1) / limit,
	})
}

func (s *UserService) Detail(c *fiber.Ctx) error {
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

func (s *UserService) Create(c *fiber.Ctx) error {
	var req model.CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	hashed, err := utils.HashPassword(req.Password)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
    }

    user := model.User{
        ID:           uuid.NewString(),
        Username:     req.Username,
        Email:        req.Email,
        FullName:     req.FullName,
        RoleID:       req.RoleID,
        PasswordHash: hashed, 
        IsActive:     true,
    }

    err = s.Repo.Create(context.Background(), &user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(user)
}

func (s *UserService) Update(c *fiber.Ctx) error {
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

func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := s.Repo.Delete(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted"})
}

func (s *UserService) UpdateRole(c *fiber.Ctx) error {
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