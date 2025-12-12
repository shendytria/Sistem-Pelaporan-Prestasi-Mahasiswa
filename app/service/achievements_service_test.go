package service

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"prestasi_mhs/app/model"
	mockRepo "prestasi_mhs/app/repository/mock"
	"prestasi_mhs/app/service/mock"
)

func setupAchievementTest() (*fiber.App, *AchievementService, *mockRepo.AchievementRepoMock, *service.StudentServiceMock) {
	repo := mockRepo.NewAchievementRepoMock()
	studentSvc := service.NewStudentServiceMock()
	svc := NewAchievementService(repo, studentSvc)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
        if c.Locals("user_id") == nil {
            c.Locals("user_id", "U1")
        }
        if c.Locals("role") == nil {
            c.Locals("role", "Mahasiswa")
        }
        c.Locals("permissions", []string{"*"}) 
        return c.Next()
    })

	app.Get("/achievement", svc.List)
	app.Get("/achievement/:id", svc.Detail)
	app.Post("/achievement", svc.Create)
	app.Put("/achievement/:id", svc.Update)
	app.Delete("/achievement/:id", svc.Delete)
	app.Put("/achievement/:id/submit", svc.Submit)
	app.Put("/achievement/:id/verify", svc.Verify)
	app.Put("/achievement/:id/reject", svc.Reject)
	app.Get("/achievement/:id/history", svc.History)
	app.Post("/achievement/:id/attachments", svc.AddAttachment)

	return app, svc, repo, studentSvc
}

func TestListAchievement(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()

	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}
	repo.ForceInsertForTest("S1")
	repo.ForceInsertForTest("S1")

	req := httptest.NewRequest("GET", "/achievement", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDetailAchievement(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()
	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}

	refID := repo.ForceInsertForTest("S1")

	req := httptest.NewRequest("GET", "/achievement/"+refID, nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCreateAchievementMahasiswa(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()

	studentSvc.Students["S1"] = model.Student{
		ID:     "S1",
		UserID: "U1",
	}

	body := model.Achievement{
		Title: "Test Achievement",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/achievement", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if len(repo.Refs) != 1 {
		t.Error("Achievement reference should be stored")
	}
}

func TestUpdateAchievementDraft(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()
	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}

	refID := repo.ForceInsertForTest("S1")
	body := `{"title":"Updated"}`
	req := httptest.NewRequest("PUT", "/achievement/"+refID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	updated := repo.Achs[refID]
	if updated.Title != "Updated" {
		t.Fatalf("expected Updated, got %s", updated.Title)
	}
}

func TestDeleteDraft(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()

	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}
	refID := repo.ForceInsertForTest("S1") 

	req := httptest.NewRequest("DELETE", "/achievement/"+refID, nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if repo.Refs[refID].Status != "deleted" {
		t.Fatal("status should be deleted")
	}
}

func TestSubmitAchievement(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()

	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}

	refID := repo.ForceInsertForTest("S1")

	req := httptest.NewRequest("PUT", "/achievement/"+refID+"/submit", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if repo.Refs[refID].Status != "submitted" {
		t.Fatal("status should change to submitted")
	}
}

func TestVerifyAchievement(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "D1")
		c.Locals("role", "Dosen Wali")
		c.Locals("permissions", []string{"*"})
		return c.Next()
	})

	studentSvc.AdvisorMapping["D1"] = []string{"S1"}
	refID := repo.ForceInsertForTestStatus("S1", "submitted")

	req := httptest.NewRequest("PUT", "/achievement/"+refID+"/verify", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if repo.Refs[refID].Status != "verified" {
		t.Fatal("status should be verified")
	}
}

func TestRejectAchievement(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "D1")
		c.Locals("role", "Dosen Wali")
		return c.Next()
	})

	studentSvc.AdvisorMapping["D1"] = []string{"S1"}
	refID := repo.ForceInsertForTestStatus("S1", "submitted")

	body := `{"reason":"Tidak valid"}`
	req := httptest.NewRequest("PUT", "/achievement/"+refID+"/reject", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if repo.Refs[refID].Status != "rejected" {
		t.Fatal("status should be rejected")
	}
}

func TestAchievementHistory(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()
	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}

	refID := repo.ForceInsertForTest("S1")

	req := httptest.NewRequest("GET", "/achievement/"+refID+"/history", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAddAttachment(t *testing.T) {
	app, _, repo, studentSvc := setupAchievementTest()
	studentSvc.Students["S1"] = model.Student{ID: "S1", UserID: "U1"}

	refID := repo.ForceInsertForTest("S1")

	body := `{"fileName":"sertifikat.pdf","fileUrl":"http://example.com/file.pdf"}`
	req := httptest.NewRequest("POST", "/achievement/"+refID+"/attachments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
