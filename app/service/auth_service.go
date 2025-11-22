package service

import (
	"context"
	"errors"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	UserSvc *UserService
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		UserSvc: NewUserService(userRepo),
	}
}

func (s *AuthService) LoginHTTP(c *fiber.Ctx) error {

	type LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	token, err := s.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {

	user, err := s.UserSvc.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}
