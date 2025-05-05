package services

import (
	"DMS/internal/dal"
	"DMS/internal/graph"
	"DMS/internal/hierarchy"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"

	"github.com/google/uuid"
)

type serviceErrorCode int

// List of error codes for methods in the services package
const (
	// The user or other entity is disabled and can't request anything
	SEIsDisabled serviceErrorCode = 0
	// The entity is exists previously
	SEExists = 1
	// A specific resource not found. (e.g. user, session, and etc)
	SENotFound = 2
	// Errors related to communicating to DBs. (e.g. connection timeout)
	SEDBError = 3
	// The user doesn't have permission to do desired action
	SENotPermission = 4
	// In some actions, the user must be ancestor to do that action but he isn't.
	SENotAncestor = 5
	// some input data are wrong
	SEWrongParameter = 6
	// Some input data are empty
	SEEmpty = 7
	// Authentication failed. e.g. the JWT is invalid or expired.
	SEAuthFailed = 8
	// The user has logged out of the session or the session has been disabled for some reason.
	SESessionExpired = 9
	// Error during encoding an entity
	SEEncodingError = 10
	// The entity deleted previously
	SEDeletedPreviously = 11
	// Claimed job position doesn't belong to the user
	SEJPNotMatchedUser = 12
	// Event not found
	SEEventNotFound = 13
	// Given job position is not owner of the event
	SEEventOwnerMismatched = 14
	// Some data in the in-memory database failed to update. Although, it doesn't effect
	// the data in the database. Means the main action is successful.
	SEInMemoryUpdateFailed = 15
	// An internal error could be database error, network error, and etc
	SEInternal = 16
	// The user/job position is denied to perform the action
	SEForbidden = 17
)

type Service struct {
	Doc           DocService
	Event         EventService
	JP            JPService
	User          UserService
	Authorization AuthorizationService
	Session       SessionService
	FilePer       FilePermissionService
}

// Create a new service
func NewService(dal *dal.DAL, hierarchy *hierarchy.HierarchyTree, cache dal.InMemoryDAL, logger l.Logger) Service {
	// Fetching job position relations
	// TODO: handle limit value in a better way
	batchSize := 500
	jpIter := dal.JP.GetJPEdgeIter(batchSize)
	changes := make(chan graph.GraphChange, batchSize)
	defer close(changes)
	go hierarchy.Graph().ProcessChanges(changes)
	for {
		val, exists := jpIter.Next()
		if !exists {
			break
		}

		responseErr := make(chan error)
		defer close(responseErr)
		changes <- graph.GraphChange{
			Type:        graph.AddEdge,
			Edge:        *jpEdge2GraphEdge(val),
			ResponseErr: responseErr,
		}
		err := <-responseErr
		if err != nil {
			logger.Panicf("Error adding edge %v: %s", *jpEdge2GraphEdge(val), err.Error())
		}
	}
	edgeCount, _ := hierarchy.Graph().Size()
	logger.Infof("Added %d vertices to the hierarchy graph", edgeCount)
	logger.Debugf("The graph:\n%s", hierarchy.Graph().String())

	session := newSSessionService(dal.Session, dal.User, logger)
	jp := newSJPService(dal.JP, hierarchy, logger)
	authorization := newSAuthorizationService(*hierarchy, dal.Permission, logger)
	event := newSEventService(dal.Event, jp, authorization, logger)
	filePermission := newSFilePermissionService(cache, session, dal.Event, authorization, logger)
	s := Service{
		Doc:           newSDocService(dal.Doc, authorization, event, jp, logger),
		Event:         event,
		JP:            jp,
		User:          newSUserService(dal.User, dal.JP, logger),
		Authorization: authorization,
		Session:       session,
		FilePer:       filePermission,
	}
	return s
}

func (s *Service) FilePermission() FilePermissionService {
	return s.FilePer
}

func jpEdge2GraphEdge(jpEdge dal.JPEdge) *graph.Edge {
	return &graph.Edge{
		Start: id2Vertex(jpEdge.Parent),
		End:   id2Vertex(jpEdge.JP),
	}
}

func id2Vertex(id m.ID) graph.Vertex {
	if id.IsNil() {
		return graph.NilVertex
	}
	return graph.Vertex(id.String())
}

// If the vertex is NilVertex, return m.NilID.
func vertex2ID(v graph.Vertex) (m.ID, error) {
	if v.Equals(graph.NilVertex) {
		return m.NilID, nil
	}
	u, err := uuid.Parse(v.String())
	if err != nil {
		return m.NilID, fmt.Errorf("invalid uuid %s: %s", v, err.Error())
	}
	id := m.ID{}
	err = id.FromUUID(u)
	if err == nil {
		return id, nil
	} else {
		return m.NilID, err
	}
}
