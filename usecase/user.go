package usecase

import (
	"andikawhy/go-user-management/repository"
	"errors"
)

func RemoveUser(deletedUserID uint64, currentUserId uint64) (*repository.User, error) {
	var userFound repository.User

	repository.DB.Where("id=?", deletedUserID).Find(&userFound)

	if userFound.ID == currentUserId {
		return nil, errors.New("cannot delete current user")
	}

	if userFound.ID == 0 {
		return nil, errors.New("user not found")
	}

	if err := repository.DB.Delete(&userFound).Error; err != nil {
		return nil, errors.New("delete user db failure")
	}

	return &userFound, nil
}

func ListUsers() (*[]repository.User, error) {
	var users []repository.User

	if err := repository.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}
