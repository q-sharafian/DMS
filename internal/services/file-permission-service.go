// This file is responsible for file related permissions, e.g. if a user wants
// to download/upload a file, the package gives the user permission to do so or not.
package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

// Specified which files are allowed to be downloaded
type allowDownload map[m.Token]bool

type allowType struct {
	FileType m.FileExtension
	IsAllow  bool
	// Maximum size of the file with with FileType in Kbytes
	MaxSize uint64
}

type FilePermissionService interface {
	// Check if each file specified in the input is allowed to be downloaded by specified
	// client that has 'AuthToken'.
	//
	// Possible error codes:
	// SEInternal- SEForbidden- SEAuthFailed
	IsAllowedDownload(accessInfo *m.DownloadReq) (allowDownload, *e.Error)

	// Check if the file type specified in the input is allowed to be uploaded and what
	// is the maximum size of each type that could be uploaded then, return the result. these details
	// are only usesable for the client with 'AuthToken' not anyone else.
	//
	// Possible error codes:
	// SEInternal- SEForbidden- SEAuthFailed
	IsAllowedUpload(accessInfo *m.UploadReq) ([]allowType, *e.Error)
}

type sFilePermissionService struct {
	cache   dal.InMemoryDAL
	session SessionService
	event   dal.EventDAL
	authz   AuthorizationService
	logger  l.Logger
}

func newSFilePermissionService(cache dal.InMemoryDAL, session SessionService, event dal.EventDAL, authzService AuthorizationService, logger l.Logger) FilePermissionService {
	return &sFilePermissionService{cache, session, event, authzService, logger}
}

// TODO: Implement cache for it
func (s *sFilePermissionService) IsAllowedDownload(accessInfo *m.DownloadReq) (allowDownload, *e.Error) {
	parsedToken, err := s.parseAuthToken(accessInfo.AuthToken)
	if err != nil {
		return nil, e.NewErrorP("failed to parse auth token: %s", SEAuthFailed, err.Error())
	}
	isAllowed, err2 := s.isAllowedAuthToken(*parsedToken)
	if err2 != nil {
		switch err2.GetCode() {
		case SEInternal, SEDBError:
			return nil, err2.SetCode(SEInternal)
		case SEAuthFailed, SENotFound:
			return nil, err2.SetCode(SEAuthFailed)
		default:
			s.logger.Warnf("Unexpected error type \"%d\" in checking auth token: %s", err2.GetCode(), err2.Error())
			return nil, err2
		}
	} else if !isAllowed {
		return nil, e.NewErrorP("the job-position %s is not allow to access event %s",
			SEForbidden, parsedToken.JobPositionID.String(), parsedToken.EventID.String())
	}

	allowDownload := make(allowDownload)
	for _, objToken := range accessInfo.ObjectTokens {
		allowDownload[objToken] = true
	}
	return allowDownload, nil
}

// TODO: Implement cache for it
func (s *sFilePermissionService) IsAllowedUpload(accessInfo *m.UploadReq) ([]allowType, *e.Error) {
	parsedToken, err := s.parseAuthToken(accessInfo.AuthToken)
	if err != nil {
		return nil, e.NewErrorP("failed to parse auth token: %s", SEAuthFailed, err.Error())
	}
	isAllowed, err2 := s.isAllowedAuthToken(*parsedToken)
	if err2 != nil {
		switch err2.GetCode() {
		case SEInternal, SEDBError:
			return nil, err2.SetCode(SEInternal)
		case SEAuthFailed, SENotFound:
			return nil, err2.SetCode(SEAuthFailed)
		default:
			s.logger.Warnf("Unexpected error type \"%d\" in checking auth token: %s", err2.GetCode(), err2.Error())
			return nil, err2
		}
	} else if !isAllowed {
		return nil, e.NewErrorP("the job-position %s is not allow to access event %s",
			SEForbidden, parsedToken.JobPositionID.String(), parsedToken.EventID.String())
	}

	var allowTypes []allowType
	for fileType, count := range accessInfo.ObjectTypes {
		allowTypes = append(allowTypes, allowType{
			FileType: fileType,
			IsAllow:  true,
			MaxSize:  uint64(count),
		})
	}
	return allowTypes, nil
}

// Check if specified job position with the given auth token exists and has access to
// the specified event.
//
// Possible error codes:
// SEAuthFailed- SEDBError- SENotFound- SEInternal
func (s *sFilePermissionService) isAllowedAuthToken(parsedAuth parsedAuthToken) (bool, *e.Error) {
	_, err := s.session.ValidateSessionJWT(parsedAuth.JWT)
	if err != nil {
		switch err.GetCode() {
		case SEAuthFailed, SEDBError:
			return false, err.AppendBegin("failed to validate auth token (error code %d)", err.GetCode())
		case SENotFound:
			return false, err.AppendBegin("it seems the session related to jwt doesn't exists")
		default:
			s.logger.Warnf("Unexpected jwt validation error code \"%d\": %s", err.GetCode(), err.Error())
			return false, err.SetCode(SEInternal)
		}
	}

	event, err2 := s.event.GetEventByID(parsedAuth.EventID)
	if err2 != nil {
		return false, e.NewErrorP("error in fetching event with id %s: %s", SEDBError, parsedAuth.EventID, err2.Error())
	} else if event == nil {
		return false, nil
	}

	isAncestor, err3 := s.authz.IsAncestor(parsedAuth.JobPositionID, event.CreatedBy)
	if err3 != nil {
		return false, err3.AppendBegin("failed to check if job-position with id %s is ancestor of %s",
			parsedAuth.JobPositionID.String(), event.CreatedBy.String()).SetCode(SEDBError)
	}
	return isAncestor, nil
}

type parsedAuthToken struct {
	JWT           m.Token
	JobPositionID m.ID
	EventID       m.ID
}

// authToken structure is as this:
// `event-id:jwt:job-position-id`. Then it must be encoded with `base64`.
func (s *sFilePermissionService) parseAuthToken(authToken m.Token) (*parsedAuthToken, error) {
	splitted, err := decodeAndSplit(authToken)
	if err != nil {
		return nil, err
	}
	parsed := parsedAuthToken{
		JWT: m.Token(splitted[1]),
	}

	id, err := m.ID{}.FromString2(splitted[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse job position id %s: %s", splitted[2], err.Error())
	}
	parsed.JobPositionID = id
	id, err = m.ID{}.FromString2(splitted[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse event id %s: %s", splitted[0], err.Error())
	}
	parsed.EventID = id

	return &parsed, nil
}

func decodeAndSplit(authToken m.Token) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(authToken.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode auth token %s with base64 decoding: %s", authToken.String(), err.Error())
	}
	decoded = bytes.TrimSpace(decoded)
	splitted := strings.SplitN(string(decoded), ":", 3)
	if len(splitted) != 3 {
		return nil, fmt.Errorf("expected 2 colon(:) character in auth token with length %d but got %d colon", len(authToken), len(splitted)-1)
	}

	return splitted, nil
}
