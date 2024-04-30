package mocks

import (
	"andikawhy/go-user-management/repository"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindByUsername(username string) repository.User {
	args := m.Called()
	return args.Get(0).(repository.User)
}

func (m *UserRepositoryMock) FindById(id uint64) repository.User {
	args := m.Called()
	return args.Get(0).(repository.User)
}

func (m *UserRepositoryMock) Delete(id uint64) repository.User {
	args := m.Called()
	return args.Get(0).(repository.User)
}

func (m *UserRepositoryMock) Save(user repository.User) repository.User {
	args := m.Called()
	return args.Get(0).(repository.User)
}

func (m *UserRepositoryMock) FindAll() []repository.User {
	args := m.Called()
	return args.Get(0).([]repository.User)
}
