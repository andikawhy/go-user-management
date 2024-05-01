package repository_test

import (
	"andikawhy/go-user-management/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
	*gorm.DB
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	m.Called(value)
	return m.DB
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	m.Called(dest, conds)
	return m.DB
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.Called(query, args)
	return m.DB
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	m.Called(value, conds)
	return m.DB
}

func TestUserRepositoryImpl_Save(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	err := db.AutoMigrate(&repository.User{})
	if err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}
	mockDB := &MockDB{DB: db}
	repo := repository.NewUserRepositoryImpl(mockDB.DB)

	user := repository.User{Username: "johndoe", Email: "john@example.com", Password: "securepassword"}

	mockDB.On("Create", &user).Return(mockDB.DB)

	assert.NotPanics(t, func() {
		repo.Save(user)
	})
}

func TestUserRepositoryImpl_Delete(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	err := db.AutoMigrate(&repository.User{})
	if err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}
	mockDB := &MockDB{DB: db}
	repo := repository.NewUserRepositoryImpl(mockDB.DB)

	mockDB.On("Where", "id=?", uint64(1)).Return(mockDB.DB)
	mockDB.On("Delete", mock.Anything).Return(mockDB.DB)

	assert.NotPanics(t, func() {
		repo.Delete(1)
	})
}

func TestUserRepositoryImpl_FindAll(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	err := db.AutoMigrate(&repository.User{})
	if err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}
	mockDB := &MockDB{DB: db}
	repo := repository.NewUserRepositoryImpl(mockDB.DB)

	mockDB.On("Find", mock.Anything).Return(mockDB.DB)

	assert.NotPanics(t, func() {
		repo.FindAll()
	})
}

func TestUserRepositoryImpl_FindById(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	err := db.AutoMigrate(&repository.User{})
	if err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}
	mockDB := &MockDB{DB: db}
	repo := repository.NewUserRepositoryImpl(mockDB.DB)

	mockDB.On("Where", "id=?", uint64(1)).Return(mockDB.DB)
	mockDB.On("Find", mock.Anything).Return(mockDB.DB)

	assert.NotPanics(t, func() {
		repo.FindById(1)
	})
}

func TestUserRepositoryImpl_FindByUsername(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	err := db.AutoMigrate(&repository.User{})
	if err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}
	mockDB := &MockDB{DB: db}
	repo := repository.NewUserRepositoryImpl(mockDB.DB)

	mockDB.On("Where", "username=?", "johndoe").Return(mockDB.DB)
	mockDB.On("Find", mock.Anything).Return(mockDB.DB)

	assert.NotPanics(t, func() {
		repo.FindByUsername("johndoe")
	})
}
