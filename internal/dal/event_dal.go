package dal

import (
	"DMS/internal/db"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"errors"
	"fmt"
)

type EventDAL interface {
	// Create event and return its id.
	CreateEvent(event *m.Event) (*m.ID, error)
	// Return some last events created by the job position id.
	// If limit be equals -1, then return all events from offset to the end.
	GetNLastEventsByJPID(jPID m.ID, limit, offset int) (*[]m.Event, error)
	GetLastApprovedEventByUserID(id m.ID) (*m.Event, *m.ApprovedEvent, error)
	// Return all created events by job position id.
	GetAllCreatedEventsByJPID(jPID m.ID) (*[]m.Event, error)
	// Return event by its id. If no error occurs and the returned event is nil, then
	// there is no corresponding event with that id.
	GetEventByID(eventID m.ID) (*m.Event, error)
	// Return some last events (specified by offset and limit)
	GetNLastEvents(limit, offset int) (*[]m.Event, error)
}

func (c cacheKey) eventByIDKey(eventID m.ID) string {
	return fmt.Sprintf("event:id:%s", eventID.String())
}

// It's an implementaion of EventDAL interface
type psqlEventDAL struct {
	db     *db.PSQLDB
	cache  *cache
	logger l.Logger
}

func newPsqlEventDAL(db *db.PSQLDB, cache *cache, logger l.Logger) *psqlEventDAL {
	return &psqlEventDAL{db, cache, logger}
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

func (d *psqlEventDAL) GetNLastEventsByJPID(jPID m.ID, limit, offset int) (*[]m.Event, error) {
	var events *[]db.Event
	result := d.db.Order("created_at desc").Limit(int(limit)).Offset(int(offset)).Where(&db.Event{
		CreatedByID: *modelID2DBID(&jPID),
	}).Find(&events)

	if result.Error != nil {
		d.logger.Debugf("Failed to get %s events for job-position-id %d (%s)", limit, jPID.String(), result.Error.Error())
		return nil, result.Error
	}
	return dbEvents2ModelEvents(events), nil
}

func (d *psqlEventDAL) GetLastApprovedEventByUserID(id m.ID) (*m.Event, *m.ApprovedEvent, error) {
	d.logger.Panicf("GetLastApprovedEventByUserID not implemented yet")
	return nil, nil, nil
}

func (d *psqlEventDAL) GetAllCreatedEventsByJPID(jpID m.ID) (*[]m.Event, error) {
	return d.GetNLastEventsByJPID(jpID, -1, 0)
}

// TODO: Test it with wrong id to know if it returns nil
func (d *psqlEventDAL) GetEventByID(eventID m.ID) (*m.Event, error) {
	var cacheKey = ck.eventByIDKey(eventID)
	var event m.Event
	if err := d.cache.read(cacheKey, &event); err != nil && !errors.Is(err, e.ErrNotFound) {
		d.logger.Debugf("Error in reading value of the key \"%s\" from the cache: %s", cacheKey, err.Error())
	} else if err == nil {
		d.logger.Debugf("Successfully read value of the key \"%s\" from the cache", cacheKey)
		return &event, nil
	}

	var dbEvent db.Event
	result := d.db.Where(&db.Event{
		BaseModel: db.BaseModel{ID: *modelID2DBID(&eventID)},
	}).Find(&dbEvent)
	if result.Error != nil {
		return nil, result.Error
	}
	event = *dbEvent2ModelEvent(&dbEvent)
	if err := d.cache.write(cacheKey, event); err != nil {
		d.logger.Debugf("Can't write an entity with key \"%s\" to cache: %s", cacheKey, err.Error())
	}
	return &event, nil
}

func (d *psqlEventDAL) GetNLastEvents(limit, offset int) (*[]m.Event, error) {
	var events *[]db.Event
	result := d.db.Order("created_at desc").Offset(offset).Limit(limit).Find(&events)
	if result.Error != nil {
		d.logger.Debugf("Failed to get %s events (%s)", limit, result.Error.Error())
		return nil, result.Error
	}
	return dbEvents2ModelEvents(events), nil
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
