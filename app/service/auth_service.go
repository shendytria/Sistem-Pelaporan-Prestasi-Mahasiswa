package service

import (
	"context"
	"github.com/google/uuid"
	"prestasi_mhs/app/repository"
	"prestasi_mhs/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	UserSvc *UserService
}

func NewAuthService(userRepo repository.UserRepo) *AuthService {
	return &AuthService{
		UserSvc: NewUserService(userRepo),
	}
}

// Login godoc
// @Summary Login user
// @Description Autentikasi user dan mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginReq true "Login request"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (s *AuthService) Login(c *fiber.Ctx) error {
	type LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ctx := context.Background()

	user, err := s.UserSvc.Repo.FindByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "user not found"})
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "invalid password"})
	}

	perms, err := s.UserSvc.Repo.GetPermissionsByRole(ctx, user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load permissions"})
	}

	roleName, err := s.UserSvc.Repo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load role name"})
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID, roleName, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	refreshToken := uuid.NewString()

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        token,
			"refreshToken": refreshToken,
			"user": fiber.Map{
				"id":          user.ID,
				"username":    user.Username,
				"fullName":    user.FullName,
				"role":        roleName,
				"permissions": perms,
			},
		},
	})
}

// Refresh Token godoc
// @Summary Refresh JWT token
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /auth/refresh [post]
func (s *AuthService) Refresh(c *fiber.Ctx) error {
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

// Get Profile godoc
// @Summary Mendapatkan profil user login
// @Security BearerAuth
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /auth/profile [get]
func (s *AuthService) Profile(c *fiber.Ctx) error {
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

// Logout godoc
// @Summary Logout user dari aplikasi
// @Security BearerAuth
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logged out"})
}
