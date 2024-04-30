package usecase

import (
	"andikawhy/go-user-management/repository"
	"errors"
	"net/http"
)

func RemoveUser(deletedUserID uint64, currentUserId uint64) (*repository.UserResponse, *repository.StandardError) {
	var userFound repository.User

	repository.DB.Where("id=?", deletedUserID).Find(&userFound)

	if userFound.ID == currentUserId {
		return nil, &repository.StandardError{Error: errors.New("cannot delete current user"), ErrorCode: http.StatusBadRequest}
	}

	if userFound.ID == 0 {
		return nil, &repository.StandardError{Error: errors.New("user not found"), ErrorCode: http.StatusBadRequest}
	}

	if err := repository.DB.Delete(&userFound).Error; err != nil {
		return nil, &repository.StandardError{Error: errors.New("delete user db failure"), ErrorCode: http.StatusInternalServerError}
	}

	userResponse := repository.UserResponse{
		ID:        userFound.ID,
		Email:     userFound.Email,
		Username:  userFound.Username,
		CreatedAt: userFound.CreatedAt,
	}

	return &userResponse, nil
}

func ListUsers() (*[]repository.UserResponse, *repository.StandardError) {
	var users []repository.User

	if err := repository.DB.Find(&users).Error; err != nil {
		return nil, &repository.StandardError{Error: err, ErrorCode: http.StatusInternalServerError}
	}

	var userResponses []repository.UserResponse
	for _, user := range users {
		userResponse := repository.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
		userResponses = append(userResponses, userResponse)
	}

	return &userResponses, nil
}
