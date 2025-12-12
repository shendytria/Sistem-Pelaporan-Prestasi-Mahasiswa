package service

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"prestasi_mhs/app/model"
	mockRepo "prestasi_mhs/app/repository/mock"
)

func setupLecturer() (*fiber.App, *LecturerService, *mockRepo.LecturerRepoMock) {
	lectRepo := mockRepo.NewLecturerRepoMock()
	studRepo := mockRepo.NewStudentRepoMock()
	svc := NewLecturerService(lectRepo, studRepo)
	app := fiber.New()

	// Middleware dummy supaya tidak kena forbidden
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Admin")
		c.Locals("permissions", []string{"*"})
		return c.Next()
	})

	app.Get("/lecturers", svc.List)
	app.Get("/lecturers/:id/advisees", svc.Advisees)

	return app, svc, lectRepo
}

func TestLecturerList(t *testing.T) {
	app, _, repo := setupLecturer()

	repo.Lecturers["L1"] = model.Lecturer{ID: "L1", LecturerID: "ABC"}
	repo.Lecturers["L2"] = model.Lecturer{ID: "L2", LecturerID: "DEF"}

	req := httptest.NewRequest("GET", "/lecturers?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestLecturerListByStudent(t *testing.T) {
	app, _, repo := setupLecturer()

	// Mock user mahasiswa
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "U999")
		c.Locals("role", "Mahasiswa")
		c.Locals("permissions", []string{"read_achievement"})
		return c.Next()
	})

	// Mock student & advisor
	repo.StudentByUser["U999"] = model.Student{AdvisorID: "L7"}
	repo.Lecturers["L7"] = model.Lecturer{ID: "L7", LecturerID: "L07"}

	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestLecturerAdvisees(t *testing.T) {
	app, _, repo := setupLecturer()

	lecturerID := uuid.NewString()
	repo.Lecturers[lecturerID] = model.Lecturer{ID: lecturerID}

	repo.Advisees[lecturerID] = []model.Advisee{
		{ID: "S1"},
		{ID: "S2"},
	}

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	data := body["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("Expected 2 advisees, got %d", len(data))
	}
}

func TestLecturerAdviseesForbidden(t *testing.T) {
	_, svc, repo := setupLecturer()

	// U1 = Dosen Wali
	appMW := fiber.New()
	appMW.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Dosen Wali")
		c.Locals("permissions", []string{"verify_achievement"})
		return c.Next()
	})
	appMW.Get("/lecturers/:id/advisees", svc.Advisees)

	// Lecturer U1
	repo.Lecturers["L1"] = model.Lecturer{ID: "L1", UserID: "U1"}

	// Try access another lecturer's advisees → must forbidden
	req := httptest.NewRequest("GET", "/lecturers/L2/advisees", nil)
	resp, _ := appMW.Test(req)

	if resp.StatusCode != 403 {
		t.Errorf("Expected 403, got %d", resp.StatusCode)
	}
}
