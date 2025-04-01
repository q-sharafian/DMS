package services

import (
	"DMS/internal/dal"
	"DMS/internal/graph"
	"DMS/internal/hierarchy"
	l "DMS/internal/logger"
	m "DMS/internal/models"
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
	JPNotMatchedUser = 12
	// Event not found
	SEEventNotFound = 13
	// Given job position is not owner of the event
	SEEventOwnerMismatched = 14
	// Some data in the in-memory database failed to update. Although, it doesn't effect
	// the data in the database. Means the main action is successful.
	InMemoryUpdateFailed = 15
)

type Service struct {
	Doc           DocService
	Event         EventService
	JP            JPService
	User          UserService
	Authorization AuthorizationService
	Session       SessionService
}

// Create a new simple service
func NewSService(dal *dal.DAL, hierarchy *hierarchy.HierarchyTree, logger l.Logger) Service {
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
	logger.Infof("Added %d edges to the graph", edgeCount)

	authorization := newSAuthorizationService(*hierarchy, dal.Permission, logger)
	jp := newSJPService(dal.JP, hierarchy, logger)
	event := newSEventService(dal.Event, jp, logger)
	s := Service{
		Doc:           newSDocService(dal.Doc, authorization, event, jp, logger),
		Event:         event,
		JP:            jp,
		User:          newSUserService(dal.User, dal.JP, logger),
		Authorization: authorization,
		Session:       newSSessionService(dal.Session, dal.User, logger),
	}
	return s
}

// If result be "NilEdge", means parent is nil. So there's not an edge.
func jpEdge2GraphEdge(jpEdge dal.JPEdge) *graph.Edge {
	if jpEdge.Parent == nil {
		return &graph.NilEdge
	}
	return &graph.Edge{
		Start: id2Vertex(*jpEdge.Parent),
		End:   id2Vertex(jpEdge.JP),
	}
}

func id2Vertex(id m.ID) graph.Vertex {
	return graph.Vertex(id.ToString())
}
