package mock

import (
	"context"
	"prestasi_mhs/app/model"
)

type LecturerRepoMock struct {
	Lecturers     map[string]model.Lecturer
	Advisees      map[string][]model.Advisee
	StudentByUser map[string]model.Student
}

func NewLecturerRepoMock() *LecturerRepoMock {
	return &LecturerRepoMock{
		Lecturers:     map[string]model.Lecturer{},
		Advisees:      map[string][]model.Advisee{},
		StudentByUser: map[string]model.Student{},
	}
}

func (m *LecturerRepoMock) FindAll(ctx context.Context) ([]model.Lecturer, error) {
	var list []model.Lecturer
	for _, l := range m.Lecturers {
		list = append(list, l)
	}
	return list, nil
}

func (m *LecturerRepoMock) FindAdvisees(ctx context.Context, lecturerID string) ([]model.Advisee, error) {
	return m.Advisees[lecturerID], nil
}

func (m *LecturerRepoMock) FindByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	for _, l := range m.Lecturers {
		if l.UserID == userID {
			return &l, nil
		}
	}
	return nil, nil
}

func (m *LecturerRepoMock) FindByID(ctx context.Context, id string) (*model.Lecturer, error) {
	l, ok := m.Lecturers[id]
	if !ok {
		return nil, nil
	}
	return &l, nil
}
