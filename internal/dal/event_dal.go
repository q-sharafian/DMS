package dal

import (
	"DMS/internal/db"
	m "DMS/internal/models"
)

type EventDAL interface {
	// Create event and return its id.
	CreateEvent(event *m.Event) (m.ID, error)
	GetLastEventByUserID(id m.ID) (*m.Event, error)
	GetLastApprovedEventByUserID(id m.ID) (*m.Event, *m.ApprovedEvent, error)
	// Return all created events by user has id.
	// limit is used to limit the number of events returned. If it be zero, return all events.
	GetAllCreatedEventsByUserID(id m.ID, limit int) (*[]m.Event, error)
}

// It's an implementaion of EventDAL interface
type psqlEventDAL struct {
	db *db.PSQLDB
}

func newPsqlEventDAL(db *db.PSQLDB) *psqlEventDAL {
	return &psqlEventDAL{db}
}

func (d *psqlEventDAL) CreateEvent(event *m.Event) (m.ID, error) {
	return m.NilID, nil
}

func (d *psqlEventDAL) GetLastEventByUserID(id m.ID) (*m.Event, error) {
	return nil, nil
}

func (d *psqlEventDAL) GetLastApprovedEventByUserID(id m.ID) (*m.Event, *m.ApprovedEvent, error) {
	return nil, nil, nil
}

func (d *psqlEventDAL) GetAllCreatedEventsByUserID(id m.ID, limit int) (*[]m.Event, error) {
	return nil, nil
}
