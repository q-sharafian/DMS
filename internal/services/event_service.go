package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"time"
)

type EventService interface {
	// Create event for specified job position and return its event ID. (In the current time)
	// Each job position could create an event.
	//
	// Possible error codes:
	// SEDBError- SENotFound
	CreateEvent(event m.Event, userID m.ID) (*m.ID, *e.Error)
	// Return job position id that created the event. He's owner of specified event.
	// If no error occurs and returned event id is nil, then there is no corresponding
	// event with this id.
	//
	// Possible error codes:
	// SEDBError
	GetEventOwner(eventID m.ID) (*m.ID, *e.Error)
	// Get some last events (according to the limit and offset values) that are
	// created by specified job position. If the job position be admin, return all events.
	// if limit be equals 0, then return all events from offset to the end.
	//
	// Possible error codes:
	// SEDBError- SEJPNotMatchedUser
	GetNLastEventsByJPID(userID, claimedJPID m.ID, limit, offset uint64) (*[]m.Event, *e.Error)
}

// It's a simple implementation of EventService interface.
// This implementation has minimum functionalities.
type sEventService struct {
	event         dal.EventDAL
	jp            JPService
	authorization AuthorizationService
	logger        l.Logger
}

// Possible error codes:
// DBError
func (s *sEventService) CreateEvent(event m.Event, userID m.ID) (*m.ID, *e.Error) {
	if isExistsUser, err := s.jp.IsExistsUserWithJP(userID, event.CreatedBy); err != nil {
		return nil, e.NewErrorP("error in checking if user exists: %s", SEDBError, err.Error())
	} else if !isExistsUser {
		return nil, e.NewErrorP("There's not any user with id %s have job position id %s",
			SENotFound, userID.String(), event.CreatedBy.String())
	}

	newEvent := m.Event{
		Name:        event.Name,
		CreatedBy:   event.CreatedBy,
		Description: event.Description,
		CreatedAt:   time.Now().UTC().Unix(),
	}
	eventID, err := s.event.CreateEvent(&newEvent)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}
	return eventID, nil
}

func (s *sEventService) GetEventOwner(eventID m.ID) (*m.ID, *e.Error) {
	event, err := s.event.GetEventByID(eventID)
	if err != nil {
		s.logger.Debugf("Failed to get event by id %s (%s)", eventID.String(), err.Error())
		return nil, e.NewErrorP(err.Error(), SEDBError)
	} else if event == nil {
		return nil, nil
	}
	return &event.CreatedBy, nil
}

func (s *sEventService) GetNLastEventsByJPID(userID, claimedJPID m.ID, limit, offset uint64) (*[]m.Event, *e.Error) {
	isAdmin, err2 := s.authorization.IsAdminJP(claimedJPID)
	if err2 != nil {
		return nil, e.NewErrorP("error in checking if the job position %s is admin: %s", SEDBError, claimedJPID.String(), err2.Error())
	}
	if isAdmin {
		s.logger.Debugf("The job position %s is admin", claimedJPID.String())
		events, err := s.event.GetNLastEvents(int(limit), int(offset))
		if err != nil {
			return nil, e.NewErrorP("failed to get some last events (limit: %d, offset: %d): %s", SEDBError, limit, offset, err.Error())
		}
		s.logger.Debugf("Got %d events. (job position is admin)", len(*events))
		return events, nil
	}

	if isExistsUser, err := s.jp.IsExistsUserWithJP(userID, claimedJPID); err != nil {
		return nil, e.NewErrorP("error in checking if user exists: %s", SEDBError, err.Error())
	} else if !isExistsUser {
		return nil, e.NewErrorP("there's not any user with id %s that have job position id %s",
			SEJPNotMatchedUser, userID.String(), claimedJPID.String())
	}

	events, err := s.event.GetNLastEventsByJPID(claimedJPID, int(limit), int(offset))
	if err != nil {
		return nil, e.NewErrorP("failed to get some last events (limit: %d, offset: %d): %s",
			SEDBError, limit, offset, err.Error())
	}
	return events, nil
}

// Create an instance of sEventService struct
func newSEventService(event dal.EventDAL, jp JPService, authz AuthorizationService, logger l.Logger) EventService {
	return &sEventService{
		event,
		jp,
		authz,
		logger,
	}
}
