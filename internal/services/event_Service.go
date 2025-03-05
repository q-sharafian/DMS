package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	m "DMS/internal/models"
	"time"
)

type EventService interface {
	// Create event for specified job position and return its event ID. (In the current time)
	CreateEvent(eventName, description string, jpID m.ID) (m.ID, *e.Error)
}

// It's a simple implementation of EventService interface.
// This implementation has minimum functionalities.
type sEventService struct {
	event dal.EventDAL
}

// Possible error codes:
// DBError
func (s *sEventService) CreateEvent(eventName, description string, jpID m.ID) (m.ID, *e.Error) {
	event := m.Event{
		Name:        eventName,
		CreatedBy:   jpID,
		Description: description,
		At:          time.Now(),
	}
	eventID, err := s.event.CreateEvent(&event)
	if err != nil {
		return m.NilID, e.NewErrorP(err.Error(), DBError)
	}
	return eventID, nil
}

// Create an instance of sEventService struct
func newsEventService(event dal.EventDAL) EventService {
	return &sEventService{
		event,
	}
}
