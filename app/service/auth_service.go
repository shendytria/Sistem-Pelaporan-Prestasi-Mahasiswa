package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
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

func (s *AuthService) Login(ctx context.Context, username, password string) (map[string]interface{}, error) {

	user, err := s.UserSvc.FindByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	perms, err := s.UserSvc.Repo.GetPermissionsByRole(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	roleName, err := s.UserSvc.Repo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID, roleName, perms)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()

	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"token":        token,
			"refreshToken": refreshToken,
			"user": map[string]interface{}{
				"id":          user.ID,
				"username":    user.Username,
				"fullName":    user.FullName,
				"role":        roleName,
				"permissions": perms,
			},
		},
	}, nil
}

func (s *AuthService) RefreshHTTP(c *fiber.Ctx) error {
	type Req struct {
		Refresh string `json:"refreshToken"`
	}
	var req Req
	c.BodyParser(&req)

	if req.Refresh == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing refreshToken"})
	}

	userID := c.Locals("user_id").(string)
	roleID := c.Locals("role_id").(string)
	roleName := c.Locals("role").(string)

	perms, err := s.UserSvc.Repo.GetPermissionsByRole(context.Background(), roleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load permissions"})
	}

	newToken, err := utils.GenerateJWT(userID, roleID, roleName, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": newToken,
	})
}

func (s *AuthService) ProfileHTTP(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roleID := c.Locals("role_id").(string)
	roleName := c.Locals("role").(string)

	user, err := s.UserSvc.Repo.FindByID(context.Background(), userID)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	perms, err := s.UserSvc.Repo.GetPermissionsByRole(context.Background(), roleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load permissions"})
	}

	return c.JSON(fiber.Map{
		"id":          user.ID,
		"username":    user.Username,
		"fullName":    user.FullName,
		"email":       user.Email,
		"role":        roleName,
		"permissions": perms,
	})
}

func (s *AuthService) LogoutHTTP(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logged out"})
}
