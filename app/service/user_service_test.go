package service

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	mockRepo "prestasi_mhs/app/repository/mock"
	"prestasi_mhs/app/model"
)

func setupUser() (*fiber.App, *UserService, *mockRepo.UserRepoMock) {
	repo := mockRepo.NewUserRepoMock()
	svc := NewUserService(repo)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Mahasiswa")
		c.Locals("permissions", []string{"*"})

		return c.Next()
	})

	app.Get("/users", svc.List)
	app.Get("/users/:id", svc.Detail)
	app.Post("/users", svc.Create)
	app.Put("/users/:id", svc.Update)
	app.Delete("/users/:id", svc.Delete)
	app.Put("/users/:id/role", svc.UpdateRole)

	return app, svc, repo
}

func TestUserList(t *testing.T) {
	app, _, repo := setupUser()

	repo.Users["U1"] = model.User{ID: "U1", Username: "john"}
	repo.Users["U2"] = model.User{ID: "U2", Username: "jane"}

	req := httptest.NewRequest("GET", "/users?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	data := body["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("Expected 2 users, got %d", len(data))
	}
}

func TestUserCreate(t *testing.T) {
	app, _, repo := setupUser()

	body := model.CreateUserReq{
		Username: "admin",
		Email:    "admin@mail.com",
		FullName: "Administrator",
		Password: "123",
		RoleID:   "R1",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if len(repo.Users) != 1 {
		t.Error("User must be inserted")
	}
}

func TestUserDetailValid(t *testing.T) {
	app, _, repo := setupUser()

	id := uuid.NewString()
	repo.Users[id] = model.User{ID: id, Username: "john"}

	req := httptest.NewRequest("GET", "/users/" + id, nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestUserDetailNotFound(t *testing.T) {
	app, _, _ := setupUser()

	req := httptest.NewRequest("GET", "/users/xxx", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}
}

func TestUserUpdate(t *testing.T) {
	app, _, repo := setupUser()

	id := uuid.NewString()
	repo.Users[id] = model.User{ID: id, Username: "old", Email: "old@mail.com"}

	updateBody := model.User{
		Username: "new",
		Email:    "new@mail.com",
		FullName: "New Name",
	}

	jsonBody, _ := json.Marshal(updateBody)
	req := httptest.NewRequest("PUT", "/users/" + id, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if repo.Users[id].Username != "new" {
		t.Errorf("Update failed, username still %s", repo.Users[id].Username)
	}
}

func TestUserDelete(t *testing.T) {
	app, _, repo := setupUser()

	id := uuid.NewString()
	repo.Users[id] = model.User{ID: id, Username: "dummy"}

	req := httptest.NewRequest("DELETE", "/users/" + id, nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if len(repo.Users) != 0 {
		t.Error("User must be deleted")
	}
}

func TestUserUpdateRole(t *testing.T) {
	app, _, repo := setupUser()

	id := uuid.NewString()
	repo.Users[id] = model.User{ID: id, Username: "john", RoleID: "R1"}

	body := `{"role_id":"R2"}`
	req := httptest.NewRequest("PUT", "/users/"+id+"/role", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if repo.Users[id].RoleID != "R2" {
		t.Errorf("Role update failed, expected R2 got %s", repo.Users[id].RoleID)
	}
}
