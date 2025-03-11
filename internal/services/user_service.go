package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type UserService interface {
	// Create a user in the hierarchy tree. If the creator user be allowed and
	// not being disabled. At the end, return the id of the created user.
	//
	// Possible error codes the function could returns:
	// IsDisabled-UserExists- DBError
	//
	// Note that users are persons that must always have created by another user.
	CreateUser(name, phone string, CreatedBy m.ID) (m.ID, *e.Error)
	// Create an andmin in the hierarchy tree. Return id of created admin.
	// Note that admins are persons that don't have created by anyperson.
	//
	// Possible error codes the function could returns:
	// IsDisabled-UserExists- DBError
	CreateAdmin(name string, phone string) (m.ID, *e.Error)
}

// It's a simple implementation of UserService interface.
// This implementation has minimum functionalities.
type sUserService struct {
	user   dal.UserDAL
	jp     dal.JPDAL
	logger l.Logger
}

// Create an admin and return the admin id. If couldn't create admin, return error.
// In this implemented method, each admin could create user
// and doesn't matter the user is allow or not.
// Note that admins are persons that don't have created by anyperson.
func (s *sUserService) CreateAdmin(name string, phone string) (m.ID, *e.Error) {
	return s.createPerson(name, phone, nil, true)
}

// Create a user and return the user id. If couldn't create user, return error. In
// this implemented method, each user could create user
// and doesn't matter the user is allow or not.
func (s *sUserService) CreateUser(name string, phone string, createdBy m.ID) (m.ID, *e.Error) {
	return s.createPerson(name, phone, &createdBy, false)
}

// Create a person and return the person id. The person could be
// a user or admin
func (s *sUserService) createPerson(name string, phone string, createdBy *m.ID, isAdmin bool) (m.ID, *e.Error) {
	// Check if there's a user with given phone-number previously.
	isExists, err := s.user.IsExistUserByPhone(phone)
	if err != nil {
		return m.NilID, (e.NewErrorP(err.Error(), SEDBError))
	}
	if isExists {
		return m.NilID, e.NewErrorP(
			"the user already exists",
			SEExists,
		)
	}

	if !isAdmin {
		isDisabled, err := s.user.IsDisabledByID(*createdBy)
		if err != nil {
			return m.NilID, (e.NewErrorP(err.Error(), SEDBError))
		}
		if isDisabled {
			return m.NilID, e.NewErrorP(
				"the user is disabled",
				SEIsDisabled,
			)
		}
	}
	newPersonID, err := s.user.CreateUser(name, phone, createdBy)
	if err != nil {
		return newPersonID, e.NewErrorP(err.Error(), SEDBError)
	}

	return newPersonID, nil
}

// Create an instance of sUserService struct
func newSUserService(user dal.UserDAL, jp dal.JPDAL, logger l.Logger) UserService {
	return &sUserService{
		user,
		jp,
		logger,
	}
}
