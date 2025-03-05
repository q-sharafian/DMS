package services

import "DMS/internal/dal"

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
func NewsService(dal *dal.DAL) Service {
	s := Service{
		Doc:   newsDocService(dal.Doc),
		Event: newsEventService(dal.Event),
		JP:    newsJPService(dal.JP),
		User:  newsUserService(dal.User, dal.JP),
	}
	return s
}
