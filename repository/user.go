package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type UserResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdat"`
}

type UserRepository interface {
	Save(user User) User
	Delete(id uint64) User
	FindById(id uint64) User
	FindByUsername(username string) User
	FindAll() []User
}

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func (t *UserRepositoryImpl) Delete(id uint64) User {
	var user User
	t.Db.Where("id=?", id).Delete(&user)
	return user
}

func (t *UserRepositoryImpl) FindAll() []User {
	var users []User
	t.Db.Find(&users)
	return users
}

func (t *UserRepositoryImpl) FindById(id uint64) User {
	var foundUser User
	t.Db.Where("id=?", id).Find(&foundUser)
	return foundUser
}

func (t *UserRepositoryImpl) FindByUsername(username string) User {
	var foundUser User
	t.Db.Where("username=?", username).Find(&foundUser)
	return foundUser
}

func (t *UserRepositoryImpl) Save(user User) User {
	t.Db.Create(&user)
	return user
}

// TODO add error handler
func NewUserRepositoryImpl(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: Db}
}
