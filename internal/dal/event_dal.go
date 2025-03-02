package dal

import (
	"DMS/internal/db"
	"time"
)

// Representation of Event entity
type Event struct {
	ID ID
	// Event name
	Name        string
	UserID      ID
	At          time.Time
	UpdatedAt   time.Time
	Description string
}

// Representation of approved event entity
type approvedEvent struct {
	EventID ID
	// Name of Person approved the event
	UserID ID
	At     time.Time
}

type EventDAL interface {
	CreateEvent(event *Event) error
	GetLastEventByUserID(id ID) (*Event, error)
	GetLastApprovedEventByUserID(id ID) (*Event, *approvedEvent, error)
	// Return all created events by user has id.
	// limit is used to limit the number of events returned. If it be zero, return all events.
	GetAllCreatedEventsByUserID(id ID, limit int) (*[]Event, error)
}

// It's an implementaion of EventDAL interface
type psqlEventDAL struct {
	db *db.PSQLDB
}

func newPsqlEventDAL(db *db.PSQLDB) *psqlEventDAL {
	return &psqlEventDAL{db}
}

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
