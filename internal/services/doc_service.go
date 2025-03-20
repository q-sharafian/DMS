package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type DocService interface {
	// Create document for specified event in the current time and return its id.
	// At this stage, just the user created the event, could create document for the event.
	//
	// Possible error codes:
	// SEDBError- SEIsDisabled
	// TODO: implement SEIsDisabled
	CreateDoc(doc *m.Doc) (*m.ID, *e.Error)
	// Return n last docs by event id iff job position id have permission to read
	// docs of the event. If eventCreatedByID be nil, we fetch event creator id from
	// the database so for better performance, it's better to pass it to avoid more
	// database query.
	//
	// Possible error codes:
	// SEDBError
	GetNLastDocByEventID(eventID m.ID, eventCreatedByID *m.ID, jpID m.ID, n int) (*[]m.Doc, *e.Error)
}

// It's a simple implementation of DocService interface.
// This implementation has minimum functionalities.
type sDocService struct {
	doc        dal.DocDAL
	logger     l.Logger
	permission PermissionService
	event      EventService
}

func (s *sDocService) CreateDoc(doc *m.Doc) (*m.ID, *e.Error) {
	eventID, err := s.doc.CreateDoc(doc)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}
	return eventID, nil
}

func (s *sDocService) GetNLastDocByEventID(eventID m.ID, eventCreatedByID *m.ID, jpID m.ID, n int) (*[]m.Doc, *e.Error) {
	if eventCreatedByID == nil {
		eventOwner, err := s.event.GetEventOwner(eventID)
		if err != nil {
			return nil, e.NewErrorP(err.Error(), SEDBError)
		} else if eventOwner == nil {
			return nil, e.NewErrorP("Event with id %s not found", SEWrongParameter, eventID.ToString())
		}
		eventCreatedByID = eventOwner
	}
	// The jpID must be the same as or an ancestor of the event creator's id.
	isAncestor := s.permission.IsAncestor(jpID, *eventCreatedByID)
	if !isAncestor {
		return nil, e.NewErrorP("You don't have permission to read docs of event with id %s",
			SENotAncestor, eventID.ToString())
	}
	docs, err := s.doc.GetNLastDocByEventID(eventID, n)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}
	return docs, nil
}

// Create an instance of sDocService struct
func newSDocService(doc dal.DocDAL, permissionService PermissionService, eventService EventService, logger l.Logger) DocService {
	return &sDocService{
		doc,
		logger,
		permissionService,
		eventService,
	}
}
