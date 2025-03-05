package services

import (
	"DMS/internal/dal"
	l "DMS/internal/logger"
)

type serviceErrorCode int

// List of error codes for methods in the services package
const (
	// The user is disabled and can't do anything
	IsDisabled serviceErrorCode = 0
	// The user is exists previously
	UserExists   = 1
	UserNotFound = 2
	// Errors related to communicating to DBs. (e.g. connection timeout)
	DBError = 3
)

type Service struct {
	Doc   DocService
	Event EventService
	JP    JPService
	User  UserService
}

// Create a new simple service
func NewsService(dal *dal.DAL, logger l.Logger) Service {
	s := Service{
		Doc:   newsDocService(dal.Doc, logger),
		Event: newsEventService(dal.Event, logger),
		JP:    newsJPService(dal.JP, logger),
		User:  newsUserService(dal.User, dal.JP, logger),
	}
	return s
}
