package service

import (
	"context"
	"prestasi_mhs/app/model"
)

type StudentServiceMock struct {
	Students        map[string]model.Student
	AdvisorMapping  map[string][]string 
}

func NewStudentServiceMock() *StudentServiceMock {
	return &StudentServiceMock{
		Students:       map[string]model.Student{},
		AdvisorMapping: map[string][]string{},
	}
}

func (m *StudentServiceMock) FindByUserID(ctx context.Context, userID string) (*model.Student, error) {
	for _, s := range m.Students {
		if s.UserID == userID {
			return &s, nil
		}
	}
	return nil, nil
}

func (m *StudentServiceMock) IsMyStudent(ctx context.Context, dosenUserID, studentID string) (bool, error) {
	list, ok := m.AdvisorMapping[dosenUserID]
	if !ok {
		return false, nil
	}
	for _, sid := range list {
		if sid == studentID {
			return true, nil
		}
	}
	return false, nil
}
