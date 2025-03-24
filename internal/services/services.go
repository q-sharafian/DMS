package services

import (
	"DMS/internal/dal"
	l "DMS/internal/logger"
	"time"
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
)

type Service struct {
	Doc        DocService
	Event      EventService
	JP         JPService
	User       UserService
	Permission PermissionService
	Session    SessionService
}

// Create a new simple service
func NewSService(dal *dal.DAL, logger l.Logger) Service {
	permission := newSPermissionService(dal.Permission, logger)
	jp := newSJPService(dal.JP, logger)
	event := newSEventService(dal.Event, jp, logger)

	s := Service{
		Doc:        newSDocService(dal.Doc, permission, event, jp, logger),
		Event:      event,
		JP:         jp,
		User:       newSUserService(dal.User, dal.JP, logger),
		Permission: permission,
		Session:    newSSessionService(dal.Session, dal.User, logger),
	}
	return s
}

// Return current unix timestamp in seconds and UTC timezone.
func getStdTime() int64 {
	return time.Now().UTC().Unix()
}
