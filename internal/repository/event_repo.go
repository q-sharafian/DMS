package repository

import (
  "DMS/internal/dal"
  "DMS/internal/models"
)

// An implementation of event operations that could be used by the application
type eventRepo struct {
  eventDAL dal.EventDAL
}

func newEventRepo(eventDAL dal.EventDAL) eventRepo {
  return eventRepo{
    eventDAL,
  }
}

// Create a new event and return its ID
func (r *eventRepo) CreateEvent(event *models.Event) (models.ID, error) {
  var err = r.eventDAL.CreateEvent(&dal.Event{
    ID:           toDALID(event.ID),
    Name:         event.Name,
    UserID:       toDALID(event.CreatedByID),
    At:           event.At,
    UpdatedAt: event.UpdatedAt,
    Description:  event.Description,
  })
  if err != nil {
    return 0, err
  }
  return event.ID, nil
}

// Return all created events by user has id.
// limit is used to limit the number of events returned. If it be zero, returns
// all events.
func (r *eventRepo) GetAllCreatedEventsByUserID(id models.ID, limit int) (*[]models.Event, error) {
  var events, err = r.eventDAL.GetAllCreatedEventsByUserID(toDALID(id), limit)
  if err != nil {
    return nil, err
  }
  var ret_events []models.Event
  for _, event := range *events {
    ret_events = append(ret_events, models.Event{
      ID:   toModelID(event.ID),
      Name: event.Name,
      // NOTE: Name of creator of the event would return as empty
      CreatedBy:    "",
      CreatedByID:  toModelID(event.UserID),
      At:           event.At,
      UpdatedAt: event.UpdatedAt,
      Description:  event.Description,
    })
  }
  return &ret_events, nil
}
