package usecase

import (
	"andikawhy/go-user-management/helper"
	"andikawhy/go-user-management/repository"
	"errors"
	"net/http"
)

type UserUsecase interface {
	RemoveUser(deletedUserID uint64, currentUserId uint64) (*repository.UserResponse, *helper.StandardError)
	ListUsers() (*[]repository.UserResponse, *helper.StandardError)
}

type UserUsecaseImpl struct {
	UserRepository repository.UserRepository
}

func (t *UserUsecaseImpl) RemoveUser(deleteUserIdRequest uint64, currentUserId uint64) (*repository.UserResponse, *helper.StandardError) {
	userFound := t.UserRepository.FindById(deleteUserIdRequest)

	if userFound.ID == currentUserId {
		return nil, &helper.StandardError{Error: errors.New("cannot delete current user"), ErrorCode: http.StatusBadRequest}
	}

	if userFound.ID == 0 {
		return nil, &helper.StandardError{Error: errors.New("user not found"), ErrorCode: http.StatusBadRequest}
	}

	t.UserRepository.Delete(deleteUserIdRequest)

	userResponse := repository.UserResponse{
		ID:        userFound.ID,
		Email:     userFound.Email,
		Username:  userFound.Username,
		CreatedAt: userFound.CreatedAt,
	}

	return &userResponse, nil
}

func (t *UserUsecaseImpl) ListUsers() (*[]repository.UserResponse, *helper.StandardError) {
	users := t.UserRepository.FindAll()

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

func NewUserUsecaseImpl(userRepository repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		UserRepository: userRepository,
	}
}
