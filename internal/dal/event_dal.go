package dal

import (
	"DMS/internal/db"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type EventDAL interface {
	// Create event and return its id.
	CreateEvent(event *m.Event) (*m.ID, error)
	// Return n last events by job position id.
	GetNLastEventsByJPID(jPID m.ID, n int) (*[]m.Event, error)
	GetLastApprovedEventByUserID(id m.ID) (*m.Event, *m.ApprovedEvent, error)
	// Return all created events by job position id.
	GetAllCreatedEventsByJPID(jPID m.ID) (*[]m.Event, error)
	// Return event by its id. If no error occurs and returned event is nil, then there
	// is no corresponding event with this id.
	GetEventByID(eventID m.ID) (*m.Event, error)
}

// It's an implementaion of EventDAL interface
type psqlEventDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlEventDAL(db *db.PSQLDB, logger l.Logger) *psqlEventDAL {
	return &psqlEventDAL{db, logger}
}

func (d *psqlEventDAL) CreateEvent(event *m.Event) (*m.ID, error) {
	newEvent := db.Event{
		Name:        event.Name,
		CreatedByID: *modelID2DBID(&event.CreatedBy),
		Description: event.Description,
	}
	result := d.db.Create(&newEvent)

	if result.Error != nil {
		d.logger.Debugf("Failed to create event for job-position-id %s (%s)", newEvent.CreatedByID.ToString(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Debugf(`It seems can't create event for user-id %s. Total rows 
				created are %d"`, newEvent.CreatedByID.ToString(), result.RowsAffected)
		return nil, e.NewSError("couldn't create event")
	}
	return dbID2ModelID(&newEvent.ID), nil
}

// If n be equals -1, then return all events
func (d *psqlEventDAL) GetNLastEventsByJPID(jPID m.ID, n int) (*[]m.Event, error) {
	var events *[]db.Event
	result := d.db.Order("created_at desc").Limit(n).Where(&db.Event{
		CreatedByID: *modelID2DBID(&jPID),
	}).Find(&events)

	if result.Error != nil {
		d.logger.Debugf("Failed to get %s events for job-position-id %d (%s)", n, jPID.ToString(), result.Error.Error())
		return nil, result.Error
	}
	return dbEvents2ModelEvents(events), nil
}

func (d *psqlEventDAL) GetLastApprovedEventByUserID(id m.ID) (*m.Event, *m.ApprovedEvent, error) {
	d.logger.Panicf("GetLastApprovedEventByUserID not implemented yet")
	return nil, nil, nil
}

func (d *psqlEventDAL) GetAllCreatedEventsByJPID(jpID m.ID) (*[]m.Event, error) {
	return d.GetNLastEventsByJPID(jpID, -1)
}

// TODO: Test it with wrong id to know if it returns nil
func (d *psqlEventDAL) GetEventByID(eventID m.ID) (*m.Event, error) {
	var event *db.Event
	result := d.db.Where(&db.Event{
		BaseModel: db.BaseModel{ID: *modelID2DBID(&eventID)},
	}).Find(event)
	if result.Error != nil {
		d.logger.Debugf("Failed to get event by id %s (%s)", eventID.ToString(), result.Error.Error())
		return nil, result.Error
	}
	return dbEvent2ModelEvent(event), nil
}

func dbEvent2ModelEvent(event *db.Event) *m.Event {
	updatedAt := event.UpdatedAt.UTC().Unix()
	return &m.Event{
		ID:          *dbID2ModelID(&event.ID),
		Name:        event.Name,
		CreatedBy:   *dbID2ModelID(&event.CreatedByID),
		Description: event.Description,
		CreatedAt:   event.CreatedAt.UTC().Unix(),
		UpdatedAt:   &updatedAt,
	}
}

func dbEvents2ModelEvents(events *[]db.Event) *[]m.Event {
	var result []m.Event
	for _, event := range *events {
		result = append(result, *dbEvent2ModelEvent(&event))
	}
	return &result
}
