package mock

import (
    "context"
    "fmt"
    "time"

    "prestasi_mhs/app/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRepoMock struct {
    Refs map[string]model.AchievementReference
    Achs map[string]model.Achievement
    Mongo map[string][]string // studentID → list of mongoIDs
	Data  map[string]model.Achievement
}

func NewAchievementRepoMock() *AchievementRepoMock {
    return &AchievementRepoMock{
        Refs: map[string]model.AchievementReference{},
        Achs: map[string]model.Achievement{},
        Mongo: map[string][]string{},
		Data:  map[string]model.Achievement{},
    }
}

func (m *AchievementRepoMock) InsertMongo(ctx context.Context, a *model.Achievement) (string, error) {
    if a.ID.IsZero() {
        a.ID = primitive.NewObjectID()
    }
    m.Achs[a.ID.Hex()] = *a
    return a.ID.Hex(), nil
}

func (m *AchievementRepoMock) FindManyMongo(ctx context.Context, ids []string) ([]model.Achievement, error) {
    var out []model.Achievement
    for _, id := range ids {
        if a, ok := m.Achs[id]; ok {
            out = append(out, a)
        }
    }
    return out, nil
}

func (m *AchievementRepoMock) FindByIDMongo(ctx context.Context, id string) (*model.Achievement, error) {
    a, ok := m.Achs[id]
    if !ok {
        return nil, nil
    }
    return &a, nil
}

func (m *AchievementRepoMock) UpdateMongo(ctx context.Context, id string, data *model.AchievementMongoUpdate) error {
    ach, ok := m.Achs[id]
    if !ok {
        return nil
    }

    if data.Title != nil {
        ach.Title = *data.Title
    }
    if data.Description != nil {
        ach.Description = *data.Description
    }

    m.Achs[id] = ach
    return nil
}

func (m *AchievementRepoMock) SoftDeleteMongo(ctx context.Context, id string) error {
    ach := m.Achs[id]
    ach.Details.Location = "deleted"
    m.Achs[id] = ach
    return nil
}

func (m *AchievementRepoMock) PushAttachmentMongo(ctx context.Context, id string, file model.AchievementFile) error {
    ach := m.Achs[id]
    ach.Attachments = append(ach.Attachments, file)
    m.Achs[id] = ach
    return nil
}

func (m *AchievementRepoMock) InsertReference(ctx context.Context, ref *model.AchievementReference) error {
    m.Refs[ref.ID] = *ref
    return nil
}

func (m *AchievementRepoMock) FindAllReferences(ctx context.Context) ([]model.AchievementReference, error) {
    var out []model.AchievementReference
    for _, r := range m.Refs {
        out = append(out, r)
    }
    return out, nil
}

func (m *AchievementRepoMock) FindReferenceByID(ctx context.Context, id string) (*model.AchievementReference, error) {
    r, ok := m.Refs[id]
    if !ok {
        return nil, nil
    }
    return &r, nil
}

func (m *AchievementRepoMock) UpdateStatus(
    ctx context.Context,
    id, status string,
    submittedAt, verifiedAt *time.Time,
    verifiedBy, note *string,
    studentID *string,
) error {
    r := m.Refs[id]
    r.Status = status
    m.Refs[id] = r
    return nil
}

func (m *AchievementRepoMock) FindMongoIDsByStudent(ctx context.Context, studentID string) ([]string, error) {
    var ids []string
    for _, r := range m.Refs {
        if r.StudentID == studentID {
            ids = append(ids, r.MongoAchievementID)
        }
    }
    return ids, nil
}

func (m *AchievementRepoMock) ForceInsertForTest(studentID string) string {
    id := fmt.Sprintf("ID-%d", len(m.Refs)+1)

    m.Refs[id] = model.AchievementReference{
        ID:                id,
        StudentID:         studentID,
        Status:            "draft",
        MongoAchievementID: id,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }

    if _, ok := m.Achs[id]; !ok {
        m.Achs[id] = model.Achievement{
            Title: "",
        }
    }

    return id
}

func (m *AchievementRepoMock) ForceInsertForTestStatus(studentID, status string) string {
    id := fmt.Sprintf("ID-%d", len(m.Refs)+1)

    m.Refs[id] = model.AchievementReference{
        ID:                id,
        StudentID:         studentID,
        Status:            status,
        MongoAchievementID: id,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }

    if _, ok := m.Achs[id]; !ok {
        m.Achs[id] = model.Achievement{}
    }

    return id
}
