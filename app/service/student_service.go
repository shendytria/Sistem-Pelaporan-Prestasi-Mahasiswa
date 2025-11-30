package service

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
    Repo     *repository.StudentRepository
    RefRepo  *repository.AchievementReferenceRepository
    MongoRepo *repository.AchievementMongoRepository
}

func NewStudentService(
    studentRepo *repository.StudentRepository,
    refRepo *repository.AchievementReferenceRepository,
    mongoRepo *repository.AchievementMongoRepository,
) *StudentService {
    return &StudentService{
        Repo:     studentRepo,
        RefRepo:  refRepo,
        MongoRepo: mongoRepo,
    }
}

func (s *StudentService) FindByUserID(ctx context.Context, userID string) (*model.Student, error) {
	return s.Repo.FindByUserID(ctx, userID)
}

func (s *StudentService) IsMyStudent(ctx context.Context, dosenUserID, studentID string) (bool, error) {
    return s.Repo.IsMyStudent(ctx, dosenUserID, studentID)
}

func (s *StudentService) FindAll(ctx context.Context) ([]model.Student, error) {
	return s.Repo.FindAll(ctx)
}

func (s *StudentService) FindByID(ctx context.Context, id string) (*model.Student, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *StudentService) ListHTTP(c *fiber.Ctx) error {
    ctx := context.Background()

    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    if page < 1 { page = 1 }
    if limit < 1 { limit = 10 }
    offset := (page - 1) * limit

    allStudents, err := s.FindAll(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    var filtered []model.Student
    if role == "Admin" {
        filtered = allStudents
    } else if role == "Dosen Wali" {
        for _, st := range allStudents {
            ok, _ := s.IsMyStudent(ctx, userID, st.ID)
            if ok {
                filtered = append(filtered, st)
            }
        }
    } else {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    total := len(filtered)

    if offset >= total {
        return c.JSON(fiber.Map{
            "data":  []model.Student{},
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
    paginated := filtered[offset:end]

    return c.JSON(fiber.Map{
        "data":  paginated,
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}

func (s *StudentService) DetailHTTP(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	roleIfc := c.Locals("role")
	userIDIfc := c.Locals("user_id")
	if roleIfc == nil || userIDIfc == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	role := roleIfc.(string)
	userID := userIDIfc.(string)

	st, err := s.FindByID(ctx, id)
	if err != nil || st == nil {
		return c.Status(404).JSON(fiber.Map{"error": "student not found"})
	}

	if role == "Admin" {
		return c.JSON(st)
	} else if role == "Dosen Wali" {
		ok, _ := s.IsMyStudent(ctx, userID, st.ID)
		if ok {
			return c.JSON(st)
		} else {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	} else {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
}

func (s *StudentService) ListByStudentHTTP(c *fiber.Ctx) error {
    ctx := context.Background()
    studentID := c.Params("id")

    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(string)

    if role == "Dosen Wali" {
        ok, err := s.IsMyStudent(ctx, userID, studentID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        if !ok {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
    }

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    if page < 1 { page = 1 }
    if limit < 1 { limit = 10 }
    offset := (page - 1) * limit

    ids, err := s.RefRepo.FindMongoIDsByStudent(ctx, studentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    total := len(ids)

    if total == 0 {
        return c.JSON(fiber.Map{
            "data":  []model.Achievement{},
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": 0,
        })
    }

    if offset >= total {
        return c.JSON(fiber.Map{
            "data":  []model.Achievement{},
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
    pageIDs := ids[offset:end]

    achs, err := s.MongoRepo.FindMany(ctx, pageIDs)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "data":  achs,
        "page":  page,
        "limit": limit,
        "total": total,
        "pages": (total + limit - 1) / limit,
    })
}

func (s *StudentService) UpdateAdvisorHTTP(c *fiber.Ctx) error {
	studentID := c.Params("id")

	type Req struct {
		AdvisorID string `json:"advisorId"`
	}
	var req Req
	c.BodyParser(&req)

	if req.AdvisorID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "advisorId required"})
	}

	err := s.Repo.UpdateAdvisor(context.Background(), studentID, req.AdvisorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "advisor updated"})
}

