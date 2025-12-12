package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"
    "prestasi_mhs/middleware"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	Repo repository.LecturerRepo
    StudentRepo repository.StudentRepo
}

func NewLecturerService(repo repository.LecturerRepo, studentRepo repository.StudentRepo) *LecturerService {
	return &LecturerService{Repo: repo, StudentRepo: studentRepo}
}

// List Lecturers godoc
// @Summary Menampilkan daftar dosen
// @Description Admin dapat melihat semua dosen, Mahasiswa hanya dapat melihat dosen pembimbingnya sendiri
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param page query int false "Nomor halaman"
// @Param limit query int false "Jumlah data per halaman"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lecturers [get]
func (s *LecturerService) List(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "read_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    switch role {

	case "Dosen Wali":
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})

	case "Mahasiswa":
		student, _ := s.StudentRepo.FindStudentByUserID(ctx, userID)
		if student == nil {
			return c.Status(404).JSON(fiber.Map{"error": "student profile not found"})
		}

		lecturer, _ := s.Repo.FindByID(ctx, student.AdvisorID)
		if lecturer == nil {
			return c.Status(404).JSON(fiber.Map{"error": "advisor not found"})
		}

		return c.JSON(fiber.Map{
			"data":  []model.Lecturer{*lecturer},
			"page":  1,
			"limit": 1,
			"total": 1,
			"pages": 1,
		})

	case "Admin":
		data, err := s.Repo.FindAll(ctx)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		offset := (page - 1) * limit
		total := len(data)

		if offset >= total {
			return c.JSON(fiber.Map{"data": []model.Lecturer{}, "page": page, "limit": limit, "total": total, "pages": (total + limit - 1) / limit})
		}

		end := offset + limit
		if end > total {
			end = total
		}

		return c.JSON(fiber.Map{
			"data":  data[offset:end],
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		})
	}

	return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
}

// List Advisees godoc
// @Summary Menampilkan daftar mahasiswa bimbingan dari dosen tertentu
// @Description Admin dapat melihat advisees dari dosen mana saja, Dosen Wali hanya dapat melihat advisees miliknya sendiri 
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Dosen"
// @Param page query int false "Nomor halaman"
// @Param limit query int false "Jumlah data"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lecturers/{id}/advisees [get]
func (s *LecturerService) Advisees(c *fiber.Ctx) error {
    if !middleware.HasPermission(c, "verify_achievement") {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    ctx := context.Background()
    lecturerID := c.Params("id")
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Dosen Wali" {
        myLecturer, err := s.Repo.FindByUserID(ctx, userID)
        if err != nil || myLecturer == nil || myLecturer.ID != lecturerID {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    advisees, err := s.Repo.FindAdvisees(ctx, lecturerID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    total := len(advisees)
    if offset >= total {
        return c.JSON(fiber.Map{
            "data":  []model.Advisee{},
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + limit - 1) / limit,
        })
    }

    end := offset + limit
    if end > total {
        end = total
    }

    return c.JSON(fiber.Map{
        "data":  advisees[offset:end],
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}