package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"
	"time"
)

type JPService interface {
	// Return all job positions the user have
	GetUserJPs(userID m.ID) ([]m.JobPosotion, *e.Error)
	// Create job position for the given user and details then, reutrn its id.
	// Two of the inputs are the ID of the (parent) JP to which the current JP belongs,
	// and another is the time the JP was created.
	CreateJP(userID, regionID, parentID m.ID, jpTitle string, createdTime time.Time, permissions m.Permission) (m.ID, *e.Error)
}

// It's a simple implementation of JPService interface.
// This implementation has minimum functionalities.
type sJPService struct {
	jp     dal.JPDAL
	logger l.Logger
}

func (s *sJPService) GetUserJPs(userID m.ID) ([]m.JobPosotion, *e.Error) {
	jps, err := s.jp.GetJPsByUserID(userID)
	if err != nil {
		return []m.JobPosotion{}, e.NewErrorP(err.Error(), DBError)
	}
	return jps, nil
}

func (s *sJPService) CreateJP(userID, regionID, parentID m.ID, jpTitle string, createdTime time.Time, permissions m.Permission) (m.ID, *e.Error) {
	jp := m.JobPosotion{
		UserID:   userID,
		ParentID: parentID,
		RegionID: regionID,
		Title:    jpTitle,
		At:       createdTime,
	}
	jpID, err := s.jp.CreateJP(&jp)
	if err != nil {
		return m.NilID, e.NewErrorFmtP(err.Error(), DBError).
			AppendBegin(
				fmt.Sprintf(
					"failed to create job position for userID %s",
					jp.UserID.ToString(),
				),
			)
	}
	return jpID, nil
}

// Create an instance of sJPService struct
func newsJPService(jp dal.JPDAL, logger l.Logger) JPService {
	return &sJPService{jp, logger}
}
