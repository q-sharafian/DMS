package services

import (
	"DMS/internal/models"
	repo "DMS/internal/repository"
)

type UserService interface {
	// Create a child for a user in the hierarchy tree. If the user be allowed and
	// not being disabled.
	CreateChildUser(user models.User) (models.User, error)
}

// It's an implementation of UserService interface.
type userService struct {
	userRepo repo.UserRepo
}

func (s *userService) CreateChildUser(user models.User) (models.User, error) {

}

// Create an instance of UserService interface
func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{userRepo}
}
