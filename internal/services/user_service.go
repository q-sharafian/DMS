package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	"DMS/internal/models"
	m "DMS/internal/models"
)

type UserService interface {
	// Create a child for a user in the hierarchy tree. If the user be allowed and
	// not being disabled.
	CreateChildUser(user models.User) (models.User, error)
}

type userServiceErrorCode int

const (
	// The user is disabled and can't do anything
	IsDisabled userServiceErrorCode = iota
)

// It's a simple implementation of UserService interface.
// It has minimum functionalities.
type sUserService struct {
	user dal.UserDAL
}

// Create a child and return it.
// In this method, each user could create child and doesn't matter the user is allow or not.
// If error code be IsDisabled, it means the user is disabled.
func (s *sUserService) CreateChildUser(user m.User) (m.User, error) {
	isDisabled, _ := s.user.IsDisabledByID(user.ID)
	if isDisabled {
		return m.User{}, e.NewError(
			"the user is disabled",
			IsDisabled,
		)
	}
	if err := s.user.CreateUser(user.Name, user.PhoneNumber); err != nil {
		return m.User{}, err
	}
	return user, nil
}

// Create an instance of UserService interface
func NewUserService(userRepo dal.UserDAL) UserService {
	return &sUserService{userRepo}
}
