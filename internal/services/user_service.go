package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"errors"
	"fmt"
)

type UserService interface {
	// Create a user. If the creator job position be allowed and
	// not being disabled. At the end, return the id of the created user.
	//
	// Possible error codes the function could returns:
	// IsDisabled-UserExists- DBError
	//
	// Note that users are persons that must always have created by another user.
	CreateUser(name string, phone m.PhoneNumber, CreatedBy m.ID) (*m.ID, *e.Error)
	// Get specified user details.
	// TODO: Add the feature that just job-positions who have permission could read user details.
	//
	// Possible error codes the function could returns:
	// DBError- SENotFound
	GetUserByID(id m.ID) (*m.User, *e.Error)
	// Create an andmin. Return id of created admin.
	// Note that admins are persons that don't have created by anyperson.
	//
	// Possible error codes the function could returns:
	// IsDisabled-UserExists- DBError
	CreateAdmin(name string, phone m.PhoneNumber) (*m.ID, *e.Error)
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
func (s *sUserService) CreateAdmin(name string, phone m.PhoneNumber) (*m.ID, *e.Error) {
	return s.createPerson(name, phone, nil, true)
}

// Create a user and return the user id. If couldn't create user, return error. In
// this implemented method, each user could create user
// and doesn't matter the user is allow or not.
func (s *sUserService) CreateUser(name string, phone m.PhoneNumber, createdBy m.ID) (*m.ID, *e.Error) {
	return s.createPerson(name, phone, &createdBy, false)
}

// Create a person and return the person id. The person could be
// a user or admin
func (s *sUserService) createPerson(name string, phone m.PhoneNumber, createdBy *m.ID, isAdmin bool) (*m.ID, *e.Error) {
	// Check if there's a user with given phone-number previously.
	isExists, err := s.user.IsExistUserByPhone(phone.ToString())
	if err != nil {
		return nil, (e.NewErrorP(err.Error(), SEDBError))
	}
	if isExists {
		return nil, e.NewErrorP(
			"the user already exists",
			SEExists,
		)
	}

	if !isAdmin {
		isDisabled, err := s.user.IsDisabledByID(*createdBy)
		if err != nil {
			return nil, (e.NewErrorP(err.Error(), SEDBError))
		}
		if isDisabled {
			return nil, e.NewErrorP(
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

func (s *sUserService) GetUserByID(id m.ID) (*m.User, *e.Error) {
	user, err := s.user.GetUserByID(id)
	if err != nil {
		if errors.Is(err, e.ErrNotFound) {
			return nil, e.NewErrorP(fmt.Sprintf("not found any user with id %s: %s", id, err.Error()), SENotFound)
		}
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}
	return user, nil
}

// Create an instance of sUserService struct
func newSUserService(user dal.UserDAL, jp dal.JPDAL, logger l.Logger) UserService {
	return &sUserService{
		user,
		jp,
		logger,
	}
}
