package service

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"prestasi_mhs/app/model"
	mockRepo "prestasi_mhs/app/repository/mock"
)

func setupStudent() (*fiber.App, *StudentService, *mockRepo.StudentRepoMock) {
	repo := mockRepo.NewStudentRepoMock()
	achMock := mockRepo.NewAchievementRepoMock()
	svc := NewStudentService(repo, achMock)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Admin")
		c.Locals("permissions", []string{"*"})
		return c.Next()
	})

	app.Get("/students", svc.List)
	app.Get("/students/:id", svc.Detail)
	app.Get("/students/:id/achievements", svc.ListByStudent)
	app.Put("/students/:id/advisor", svc.UpdateAdvisor)

	return app, svc, repo
}

func TestStudentList(t *testing.T) {
	app, _, repo := setupStudent()

	repo.Students["S1"] = model.Student{ID: "S1", StudentID: "123"}
	repo.Students["S2"] = model.Student{ID: "S2", StudentID: "456"}

	req := httptest.NewRequest("GET", "/students?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	data := body["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("Expected 2 students, got %d", len(data))
	}
}

func TestStudentDetailFound(t *testing.T) {
	app, _, repo := setupStudent()

	id := uuid.NewString()
	repo.Students[id] = model.Student{ID: id, StudentID: "999"}

	req := httptest.NewRequest("GET", "/students/"+id, nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestStudentDetailNotFound(t *testing.T) {
	app, _, _ := setupStudent()

	req := httptest.NewRequest("GET", "/students/xyz", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}
}

func TestStudentListByStudent(t *testing.T) {
	app, svc, repo := setupStudent()

	repo.Students["S1"] = model.Student{ID: "S1"}

	achMock := svc.AchievementRepo.(*mockRepo.AchievementRepoMock)
	achMock.ForceInsertForTest("S1")
	achMock.ForceInsertForTest("S1")

	req := httptest.NewRequest("GET", "/students/S1/achievements?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestStudentUpdateAdvisor(t *testing.T) {
	app, _, repo := setupStudent()

	id := uuid.NewString()
	repo.Students[id] = model.Student{ID: id}

	body := `{"advisorId":"L001"}`
	req := httptest.NewRequest("PUT", "/students/"+id+"/advisor", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if repo.Students[id].AdvisorID != "L001" {
		t.Errorf("Advisor update failed, expected L001 got %s", repo.Students[id].AdvisorID)
	}
}
