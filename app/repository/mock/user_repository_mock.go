package mock

import (
	"context"
	"prestasi_mhs/app/model"
)

type UserRepoMock struct {
	Users             map[string]model.User
	PermissionsByRole map[string][]string 
	RoleNameByID      map[string]string   
}

func NewUserRepoMock() *UserRepoMock {
	return &UserRepoMock{
		Users:             make(map[string]model.User),
		PermissionsByRole: make(map[string][]string),
		RoleNameByID:      make(map[string]string),
	}
}

func (m *UserRepoMock) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	for _, u := range m.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, nil
}

func (m *UserRepoMock) FindAll(ctx context.Context, limit, offset int) ([]model.User, int, error) {
	var list []model.User
	for _, u := range m.Users {
		list = append(list, u)
	}
	return list, len(list), nil
}

func (m *UserRepoMock) FindByID(ctx context.Context, id string) (*model.User, error) {
	if u, ok := m.Users[id]; ok {
		return &u, nil
	}
	return nil, nil
}

func (m *UserRepoMock) Create(ctx context.Context, u *model.User) error {
	m.Users[u.ID] = *u
	return nil
}

func (m *UserRepoMock) Update(ctx context.Context, u *model.User) error {
	m.Users[u.ID] = *u
	return nil
}

func (m *UserRepoMock) Delete(ctx context.Context, id string) error {
	delete(m.Users, id)
	return nil
}

func (m *UserRepoMock) UpdateRole(ctx context.Context, id string, roleID string) error {
	u := m.Users[id]
	u.RoleID = roleID
	m.Users[id] = u
	return nil
}

func (m *UserRepoMock) GetPermissionsByRole(ctx context.Context, roleID string) ([]string, error) {
	return m.PermissionsByRole[roleID], nil
}

func (m *UserRepoMock) GetRoleName(ctx context.Context, roleID string) (string, error) {
	return m.RoleNameByID[roleID], nil
}
