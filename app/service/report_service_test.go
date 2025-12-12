package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	mockRepo "prestasi_mhs/app/repository/mock"
	"prestasi_mhs/app/model"
)

func setupReport() (*ReportService, *mockRepo.ReportRepoMock, *mockRepo.StudentRepoMock, *mockRepo.LecturerRepoMock) {
	reportRepo := mockRepo.NewReportRepoMock()
	studentRepo := mockRepo.NewStudentRepoMock()
	lecturerRepo := mockRepo.NewLecturerRepoMock()

	svc := NewReportService(reportRepo, studentRepo, lecturerRepo)
	return svc, reportRepo, studentRepo, lecturerRepo
}

func TestStatistics_Admin(t *testing.T) {
	svc, reportRepo, _, _ := setupReport()
	app := fiber.New()

	reportRepo.MockStatistics = &model.StatisticResponse{TotalAchievements: 99}

	app.Get("/reports/statistics", func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Admin")
		c.Locals("permissions", []string{"read_achievement"})
		return svc.Statistics(c)
	})

	req := httptest.NewRequest("GET", "/reports/statistics", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestStatistics_DosenWali(t *testing.T) {
    svc, reportRepo, _, lecturerRepo := setupReport()
    app := fiber.New()

    lecturerRepo.Lecturers["L1"] = model.Lecturer{
        ID:     "L1",
        UserID: "U1",
    }

    reportRepo.MockStatistics = &model.StatisticResponse{Verified: 3}

    app.Get("/reports/statistics", func(c *fiber.Ctx) error {
        c.Locals("user_id", "U1")
        c.Locals("role", "Dosen Wali")
        c.Locals("permissions", []string{"read_achievement"})
        return svc.Statistics(c)
    })

    req := httptest.NewRequest("GET", "/reports/statistics", nil)
    resp, _ := app.Test(req)

    if resp.StatusCode != 200 {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
}

func TestStatistics_Mahasiswa(t *testing.T) {
	svc, reportRepo, studentRepo, _ := setupReport()
	app := fiber.New()

	studentRepo.StudentsByUser["U1"] = model.Student{ID: "S01"}
	reportRepo.MockStatistics = &model.StatisticResponse{Rejected: 2}

	app.Get("/reports/statistics", func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Mahasiswa")
		c.Locals("permissions", []string{"read_achievement"})
		return svc.Statistics(c)
	})

	req := httptest.NewRequest("GET", "/reports/statistics", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestStudentReport_Admin(t *testing.T) {
	svc, reportRepo, _, _ := setupReport()
	app := fiber.New()

	reportRepo.MockStudentReport = &model.StudentReportResponse{StudentID: "S1"}

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Admin")
		c.Locals("permissions", []string{"read_achievement"})
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest("GET", "/reports/student/S1", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestStudentReport_ForbiddenForMahasiswa(t *testing.T) {
	svc, _, studentRepo, _ := setupReport()
	app := fiber.New()

	studentRepo.StudentsByUser["U1"] = model.Student{ID: "S99"} 

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Mahasiswa")
		c.Locals("permissions", []string{"read_achievement"})
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest("GET", "/reports/student/S1", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 403 {
		t.Errorf("Expected 403, got %d", resp.StatusCode)
	}
}

func TestStudentReport_ForbiddenForDosenWali(t *testing.T) {
	svc, _, studentRepo, lecturerRepo := setupReport()
	app := fiber.New()

	lecturerRepo.Lecturers["L10"] = model.Lecturer{ID: "L10", UserID: "U1"}
	studentRepo.Students["S1"] = model.Student{ID: "S1", AdvisorID: "L20"} 

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "U1")
		c.Locals("role", "Dosen Wali")
		c.Locals("permissions", []string{"read_achievement"})
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest("GET", "/reports/student/S1", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 403 {
		t.Errorf("Expected 403, got %d", resp.StatusCode)
	}
}
