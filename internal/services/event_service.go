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
	// SEDBError
	CreateEvent(event m.Event) (*m.ID, *e.Error)
	// Return job position id that created the event. He's owner of specified event.
	// If no error occurs and returned event id is nil, then there is no corresponding
	// event with this id.
	//
	// Possible error codes:
	// SEDBError
	GetEventOwner(eventID m.ID) (*m.ID, *e.Error)
}

// It's a simple implementation of EventService interface.
// This implementation has minimum functionalities.
type sEventService struct {
	event  dal.EventDAL
	logger l.Logger
}

// Possible error codes:
// DBError
func (s *sEventService) CreateEvent(event m.Event) (*m.ID, *e.Error) {
	newEvent := m.Event{
		Name:        event.Name,
		CreatedBy:   event.CreatedBy,
		Description: event.Description,
		CreatedAt:   time.Now(),
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
		s.logger.Debugf("Failed to get event by id %s (%s)", eventID.ToString(), err.Error())
		return nil, e.NewErrorP(err.Error(), SEDBError)
	} else if event == nil {
		return nil, nil
	}
	return &event.CreatedBy, nil
}

// Create an instance of sEventService struct
func newSEventService(event dal.EventDAL, logger l.Logger) EventService {
	return &sEventService{
		event,
		logger,
	}
}
