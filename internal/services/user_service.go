package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type UserService interface {
	// Create a user in the hierarchy tree. If the creator user be allowed and
	// not being disabled.
	CreateUser(name, phone string, CreatedBy m.ID) (m.ID, *e.Error)
}

// It's a simple implementation of UserService interface.
// This implementation has minimum functionalities.
type sUserService struct {
	user   dal.UserDAL
	jp     dal.JPDAL
	logger *l.Logger
}

// Create a user and return its id. If couldn't create user, return error.
// In this implemented method, each user could create user and doesn't matter the
// user is allow or not.
// Possible error codes the function could return:
// IsDisabled-UserExists- DBError
func (s sUserService) CreateUser(name string, phone string, CreatedBy m.ID) (m.ID, *e.Error) {
	// Check if there's a user with given phone-number previously.
	isExists, err := s.user.IsExistUserByPhone(phone)
	if err != nil {
		return m.NilID, (e.NewErrorP(err.Error(), DBError))
	}
	if isExists {
		return m.NilID, e.NewErrorP(
			"the user already exists",
			UserExists,
		)
	}

	isDisabled, _ := s.user.IsDisabledByID(CreatedBy)
	if isDisabled {
		return m.NilID, e.NewErrorP(
			"the user is disabled",
			IsDisabled,
		)
	}
	newUserID, err := s.user.CreateUser(name, phone, CreatedBy)
	if err != nil {
		return newUserID, e.NewErrorP(err.Error(), DBError)
	}

	return newUserID, nil
}

// Create an instance of sUserService struct
func newsUserService(user dal.UserDAL, jp dal.JPDAL, logger *l.Logger) UserService {
	return &sUserService{
		user,
		jp,
		logger,
	}
}
