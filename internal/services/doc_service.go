package services

import (
	"DMS/internal/common"
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type DocService interface {
	// Create document for specified event and job position in the current time and return its id.
	// At this stage, just the user created the event, could create document for the event.
	//
	// Possible error codes:
	// SEDBError- SEIsDisabled- SEEventOwnerMismatched- SENotFound
	// TODO: implement SEIsDisabled
	CreateDoc(doc *m.Doc, userID m.ID) (*m.ID, *e.Error)
	// Return n last docs by event id iff job position id have permission to read
	// docs of the event. If eventCreatedByID be nil, we fetch event creator id from
	// the database so for better performance, it's better to pass it to avoid more
	// database query. jpID is a job position id that belongs to the userID.
	//
	// Possible error codes:
	// SEDBError- SEJPNotMatchedUser- SENotAncestor
	GetNLastDocByEventID(eventID, userID m.ID, eventCreatedByID *m.ID, jpID m.ID, n int) (*[]m.Doc, *e.Error)
	// Get some last documents (according to the limit and offset values) that are
	// accessible for the job position. If the job position be admin, he accesses to all docs.
	//
	// Possible error codes:
	// SEDBError- SEJPNotMatchedUser
	GetNLastDocs(userID, claimedJPID m.ID, limit, offset uint64) (*[]m.DocWithSomeDetails, *e.Error)
}

// It's a simple implementation of DocService interface.
// This implementation has minimum functionalities.
type sDocService struct {
	doc           dal.DocDAL
	logger        l.Logger
	authorization AuthorizationService
	event         EventService
	jp            JPService
}

func (s *sDocService) CreateDoc(doc *m.Doc, userID m.ID) (*m.ID, *e.Error) {
	if isExistsUser, err := s.jp.IsExistsUserWithJP(userID, doc.CreatedBy); err != nil {
		return nil, e.NewErrorP("error in checking if user exists: %s", SEDBError, err.Error())
	} else if !isExistsUser {
		return nil, e.NewErrorP("there's not any user with id %s that have job position id %s",
			SENotFound, userID.String(), doc.CreatedBy.String())
	}
	// Just job position who created the event could create document for that.
	if eventOwner, err := s.event.GetEventOwner(doc.EventID); err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	} else if eventOwner == nil {
		return nil, e.NewErrorP("Event with id %s not found", SEEventNotFound, doc.EventID.String())
	} else if *eventOwner != doc.CreatedBy {
		return nil, e.NewErrorP("The job position %s is not ownwe of the event %s",
			SEEventOwnerMismatched, doc.CreatedBy.String(), doc.EventID.String())
	}

	eventID, err := s.doc.CreateDoc(doc)
	if err != nil {
		return nil, e.NewErrorP("failed to create doc: %s", SEDBError, err.Error())
	}
	return eventID, nil
}

func (s *sDocService) GetNLastDocByEventID(eventID, userID m.ID, eventCreatedByID *m.ID, jpID m.ID, n int) (*[]m.Doc, *e.Error) {
	if eventCreatedByID == nil {
		eventOwner, err := s.event.GetEventOwner(eventID)
		if err != nil {
			return nil, e.NewErrorP(err.Error(), SEDBError)
		} else if eventOwner == nil {
			return nil, e.NewErrorP("event with id %s not found", SEWrongParameter, eventID.String())
		}
		eventCreatedByID = eventOwner
	}

	// Validate if claimed jon position id belongs to the specified user.
	// isExistsUser, err := s.jp.IsExistsUserWithJP(userID, jpID)
	// if err != nil {
	// 	return nil, e.NewErrorP("error in checking if user exists: %s", SEDBError, err.Error())
	// } else if !isExistsUser {
	// 	return nil, e.NewErrorP("there's not any user with id %s have job position id %s",
	// 		JPNotMatchedUser, jpID.ToString(), eventCreatedByID.ToString())
	// }

	// The jpID must be the same as or an ancestor of the event creator's id.
	if isAncestor, err2 := s.authorization.IsAncestor(jpID, *eventCreatedByID); err2 != nil {
		return nil, e.NewErrorP(err2.Error(), SEDBError)
	} else if !isAncestor {
		return nil, e.NewErrorP("you don't have permission to read docs of event with id %s",
			SENotAncestor, eventID.String())
	}
	docs, err := s.doc.GetNLastDocByEventID(eventID, n)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}
	return docs, nil
}

func (s *sDocService) GetNLastDocs(userID, claimedJPID m.ID, limit, offset uint64) (*[]m.DocWithSomeDetails, *e.Error) {
	if isExistsUser, err := s.jp.IsExistsUserWithJP(userID, claimedJPID); err != nil {
		return nil, e.NewErrorP("error in checking if user exists: %s", SEDBError, err.Error())
	} else if !isExistsUser {
		return nil, e.NewErrorP("there's not any user with id %s that have job position id %s",
			SEJPNotMatchedUser, userID.String(), claimedJPID.String())
	}

	isAdmin, err2 := s.authorization.IsAdminJP(claimedJPID)
	if err2 != nil {
		return nil, e.NewErrorP("error in checking if the job position %s is admin: %s", SEDBError, claimedJPID.String(), err2.Error())
	}
	if isAdmin {
		s.logger.Debugf("The job position %s is admin", claimedJPID.String())
		docs, err := s.doc.GetNLastDocs(int(limit), int(offset))
		if err != nil {
			return nil, e.NewErrorP("failed to get some last docs (limit: %d, offset: %d): %s", SEDBError, limit, offset, err.Error())
		}
		s.logger.Debugf("Got %d docs. (job position is admin)", len(*docs))
		return docs, nil
	}

	childJPIDs, err2 := s.authorization.GetNestedChilds(claimedJPID)
	if err2 != nil {
		return nil, e.NewErrorP("can't fetch nested childs: %s", SEDBError, err2.Error())
	} else {
		s.logger.Debugf("Fetched %d child job positions", len(childJPIDs))
		maxToShow := min(len(childJPIDs), 10)
		s.logger.Debugf("%d nested job positions for specified job position %s: %+v",
			maxToShow, claimedJPID, childJPIDs[:maxToShow])
	}

	docGroups := make([]*[]m.DocWithSomeDetails, 0)
	for _, childJPID := range childJPIDs {
		// TODO: Optimize and edit the number of documents to fetch from each job position id
		// such that improve performance.
		docs, err := s.doc.GetNLastDocsByJPID(childJPID, int(offset), int(limit))
		if err != nil {
			return nil, e.NewErrorP("failed to get %s docs with offset %s for child job position: %s",
				SEDBError, limit, offset, err.Error())
		}
		docGroups = append(docGroups, docs)
	}

	orderedDocs := common.MergeOrderedGroups(&docGroups, func(a, b m.DocWithSomeDetails) bool {
		return a.CreatedAt < b.CreatedAt
	})
	if len(*orderedDocs) < int(limit) {
		return orderedDocs, nil
	} else {
		limitedDocs := (*orderedDocs)[:limit]
		return &limitedDocs, nil
	}
}

// Create an instance of sDocService struct
func newSDocService(doc dal.DocDAL, permissionService AuthorizationService, eventService EventService,
	jpService JPService, logger l.Logger) DocService {
	return &sDocService{
		doc,
		logger,
		permissionService,
		eventService,
		jpService,
	}
}
