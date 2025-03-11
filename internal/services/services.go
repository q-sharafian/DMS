package services

import (
	"DMS/internal/dal"
	l "DMS/internal/logger"
)

type serviceErrorCode int

// List of error codes for methods in the services package
const (
	// The user is disabled and can't request anything
	SEIsDisabled serviceErrorCode = 0
	// The entity is exists previously
	SEExists   = 1
	SENotFound = 2
	// Errors related to communicating to DBs. (e.g. connection timeout)
	SEDBError = 3
	// The user doesn't have permission to do desired action
	SENotPermission = 4
	// In some actions, the user must be ancestor to do that action but he isn't.
	SENotAncestor = 5
	// some input data are wrong
	SEWrongParameter = 6
)

type Service struct {
	Doc        DocService
	Event      EventService
	JP         JPService
	User       UserService
	Permission PermissionService
}

// Create a new simple service
func NewSService(dal *dal.DAL, logger l.Logger) Service {
	permission := newSPermissionService(dal.Permission, logger)
	event := newSEventService(dal.Event, logger)

	s := Service{
		Doc:        newSDocService(dal.Doc, permission, event, logger),
		Event:      event,
		JP:         newSJPService(dal.JP, logger),
		User:       newSUserService(dal.User, dal.JP, logger),
		Permission: permission,
	}
	return s
}
