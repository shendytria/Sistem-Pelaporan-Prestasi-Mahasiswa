package mock

import (
	"context"
	"prestasi_mhs/app/model"
)

type StudentRepoMock struct {
	Students map[string]model.Student
	StudentsByUser map[string]model.Student 
	Advisors       map[string]string 
}

func NewStudentRepoMock() *StudentRepoMock {
	return &StudentRepoMock{
		Students:      map[string]model.Student{},
		StudentsByUser: map[string]model.Student{}, 
		Advisors:       map[string]string{},
	}
}

func (m *StudentRepoMock) FindByUserID(ctx context.Context, userID string) (*model.Student, error) {
    if s, ok := m.StudentsByUser[userID]; ok {
        return &s, nil
    }
    return nil, nil
}

func (m *StudentRepoMock) IsMyStudent(ctx context.Context, dosenID, studentID string) (bool, error) {
	s, ok := m.Students[studentID]
	if !ok {
		return false, nil
	}
	return s.AdvisorID == dosenID, nil
}

func (m *StudentRepoMock) FindAll(ctx context.Context) ([]model.Student, error) {
	var list []model.Student
	for _, s := range m.Students {
		list = append(list, s)
	}
	return list, nil
}

func (m *StudentRepoMock) FindByID(ctx context.Context, id string) (*model.Student, error) {
	s, ok := m.Students[id]
	if !ok {
		return nil, nil
	}
	return &s, nil
}

func (m *StudentRepoMock) UpdateAdvisor(ctx context.Context, studentID, lecturerID string) error {
	s := m.Students[studentID]
	s.AdvisorID = lecturerID
	m.Students[studentID] = s
	return nil
}

func (m *StudentRepoMock) FindStudentByUserID(ctx context.Context, userID string) (*model.Student, error) {
    if s, ok := m.StudentsByUser[userID]; ok {
        return &s, nil
    }
    return nil, nil
}

func (m *StudentRepoMock) CheckAdvisor(ctx context.Context, lecturerID, studentID string) (bool, error) {
    st, ok := m.Students[studentID]
    if !ok {
        return false, nil
    }
    return st.AdvisorID == lecturerID, nil
}
