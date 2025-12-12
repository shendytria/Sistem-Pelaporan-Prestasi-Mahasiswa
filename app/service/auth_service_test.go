package service

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"prestasi_mhs/app/model"
	mockRepo "prestasi_mhs/app/repository/mock"
	"prestasi_mhs/utils"
)

func setupAuth() (*fiber.App, *AuthService, *mockRepo.UserRepoMock) {
    repo := mockRepo.NewUserRepoMock()
    svc := NewAuthService(repo)

    app := fiber.New()

    app.Use(func(c *fiber.Ctx) error {
        if c.Locals("user_id") == nil {
            c.Locals("user_id", "U1")
        }
        if c.Locals("role_id") == nil {
            c.Locals("role_id", "R1")
        }
        if c.Locals("role") == nil {
            c.Locals("role", "Admin")
        }
        if c.Locals("permissions") == nil {
            c.Locals("permissions", []string{"*"})
        }
        return c.Next()
    })

    app.Post("/login", svc.Login)
    app.Post("/auth/refresh", svc.Refresh)
    app.Get("/auth/profile", svc.Profile)
    app.Post("/auth/logout", svc.Logout)

    return app, svc, repo
}

func TestLoginSuccess(t *testing.T) {
	app, _, repo := setupAuth()

	hash, _ := utils.HashPassword("123")

	repo.Users["U1"] = model.User{
		ID:           "U1",
		Username:     "admin",
		PasswordHash: hash,
		FullName:     "Administrator",
		RoleID:       "R1",
	}

	repo.PermissionsByRole["R1"] = []string{"user:read", "user:create"}
	repo.RoleNameByID["R1"] = "Admin"

	body := map[string]string{
		"username": "admin",
		"password": "123",
	}
	j, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	app, _, repo := setupAuth()

	hash, _ := utils.HashPassword("123")
	repo.Users["U1"] = model.User{
		ID:           "U1",
		Username:     "admin",
		PasswordHash: hash,
	}

	body := map[string]string{
		"username": "admin",
		"password": "salah",
	}
	j, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("Expected 401, got %d", resp.StatusCode)
	}
}

func TestLoginUserNotFound(t *testing.T) {
	app, _, _ := setupAuth()

	body := map[string]string{
		"username": "unknown",
		"password": "123",
	}
	j, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("Expected 401, got %d", resp.StatusCode)
	}
}

func TestRefreshSuccess(t *testing.T) {
    app, _, repo := setupAuth()

    repo.PermissionsByRole["R1"] = []string{"user:read"}
    repo.RoleNameByID["R1"] = "Admin"

    body := map[string]string{"refreshToken": "dummy"}
    b, _ := json.Marshal(body)

    req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json")

    resp, _ := app.Test(req)

    if resp.StatusCode != 200 {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
}

func TestProfile(t *testing.T) {
	app, _, repo := setupAuth()

	repo.Users["U1"] = model.User{ID: "U1", Username: "admin", FullName: "Administrator", Email: "admin@mail.com"}
	repo.PermissionsByRole["R1"] = []string{"user:read"}
	repo.RoleNameByID["R1"] = "Admin"

	req := httptest.NewRequest("GET", "/auth/profile", nil)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role_id", "R1")
		c.Locals("role", "Admin")
		return c.Next()
	})

	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	app, _, _ := setupAuth()

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
