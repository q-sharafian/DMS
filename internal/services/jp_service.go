package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"
)

type JPService interface {
	// Return all job positions the user have
	//
	// Possible error codes the function could returns:
	// SEDBError- SENotFound
	GetUserJPs(user *m.User) (*[]m.UserJobPosition, *e.Error)
	// Create user job position with its permissions for the given user and details then, reutrn its id.
	//
	// Possible error codes the function could returns:
	// SEDBError- SENotFound
	CreateUserJP(jp *m.UserJobPosition, permissions *m.Permission) (*m.ID, *e.Error)
	// Create admin job position with its permissions for the given user and details then, reutrn its id.
	//
	// Possible error codes the function could returns:
	// SEDBError
	CreateAdminJP(jp *m.AdminJobPosition, permissions *m.Permission) (*m.ID, *e.Error)
	// Return true if a job position with given ID belongs to a user with given ID.
	//
	// Possible error codes:
	// SEDBError
	IsExistsUserWithJP(userID, jpID m.ID) (bool, error)
}

// It's a simple implementation of JPService interface.
// This implementation has minimum functionalities.
type sJPService struct {
	jp     dal.JPDAL
	logger l.Logger
}

func (s *sJPService) GetUserJPs(user *m.User) (*[]m.UserJobPosition, *e.Error) {
	jps, err := s.jp.GetJPsByUser(user)
	if err != nil {
		return &[]m.UserJobPosition{}, e.NewErrorP(err.Error(), SEDBError)
	} else if jps == nil {
		return jps, e.NewErrorP("no job positions found for user %+v", SENotFound, user)
	}
	return jps, nil
}

// Note that in this implementation, createdTime value doesn't matter and createdTime
// is always the current time.
func (s *sJPService) CreateUserJP(jp *m.UserJobPosition, permissions *m.Permission) (*m.ID, *e.Error) {
	jpID, err := s.jp.CreateUserJPWithPermissions(jp, permissions)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError).
			AppendBegin(
				fmt.Sprintf(
					"failed to create job position for userID %s",
					jp.UserID.ToString(),
				),
			)
	}
	return jpID, nil
}

// Note that in this implementation, createdTime value doesn't matter and createdTime
// is always the current time.
func (s *sJPService) CreateAdminJP(jp *m.AdminJobPosition, permissions *m.Permission) (*m.ID, *e.Error) {
	jpID, err := s.jp.CreateAdminJPWithPermissions(jp, permissions)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError).
			AppendBegin(
				fmt.Sprintf(
					"failed to create job position for userID %s",
					jp.UserID.ToString(),
				),
			)
	}
	return jpID, nil
}

func (s *sJPService) IsExistsUserWithJP(userID, jpID m.ID) (bool, error) {
	isExists, err := s.jp.IsExistsUserWithJP(userID, jpID)
	if err != nil {
		return false, e.NewErrorP(err.Error(), SEDBError)
	}
	return isExists, nil
}

// Create an instance of sJPService struct
func newSJPService(jp dal.JPDAL, logger l.Logger) JPService {
	return &sJPService{jp, logger}
}
