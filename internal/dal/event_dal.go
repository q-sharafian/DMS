package dal

import (
  "time"
)

// Representation of Event entity
type Event struct {
  ID ID `db:"id"`
  // Event name
  Name         string    `db:"name"`
  UserID       ID        `db:"user_id"`
  At           time.Time `db:"at"`
  UpdatedAt time.Time `db:"last_change_at"`
  Description  string    `db:"description"`
}

// Representation of approved event entity
type approvedEvent struct {
  EventID ID `db:"event_id"`
  // Name of Person approved the event
  UserID ID        `db:"user_id"`
  At     time.Time `db:"at"`
}

type EventDAL interface {
  CreateEvent(event *Event) error
  GetLastEventByUserID(id ID) (*Event, error)
  GetLastApprovedEventByUserID(id ID) (*Event, *approvedEvent, error)
  // Return all created events by user has id.
  // limit is used to limit the number of events returned. If it be zero, return all events.
  GetAllCreatedEventsByUserID(id ID, limit int) (*[]Event, error)
}

type psqlEventDAL struct{}

func (d *psqlEventDAL) CreateEvent(event *Event) error {
  return nil
}

func (d *psqlEventDAL) GetLastEventByUserID(id ID) (*Event, error) {
  return nil, nil
}

func (d *psqlEventDAL) GetLastApprovedEventByUserID(id ID) (*Event, *approvedEvent, error) {
  return nil, nil, nil
}

func (d *psqlEventDAL) GetAllCreatedEventsByUserID(id ID, limit int) (*[]Event, error) {
  return nil, nil
}
