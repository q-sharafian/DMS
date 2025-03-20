package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Contains interface for all functionalities related to sessions.
type SessionService interface {
	// Create a login for the specified user and return its id.
	//
	// Possible error codes:
	// SEDBError- SENotFound
	CreateSessionJustByPhone(session *m.Session, phone m.PhoneNumber) (*m.ID, *e.Error)
	// Delete the session associated with the JWT. Note that the user id that sends the session deletion
	// The request must match the user id that the session is created for.
	//
	// Possible error codes:
	// SEDBError- SENotFound
	DeleteSession(jwt *m.JWT) *e.Error
	// Validate session based on the input jwt token. We must remove any prefix like "Bearer " from the
	// input JWT token before calling the method wih that value.
	//
	// Possible error codes:
	// SEAuthFailed- SENotFound
	ValidateSessionJWT(token string) (*m.JWT, *e.Error)
	// If both error and session be nil, means there's not any matched session.
	// (whether disabled, removed, and etc.)
	//
	// Possible error codes:
	// SEDBError
	GetSessionByID(sessionID *m.ID) (*m.Session, *e.Error)
	// GetSessionByToken(token string) (*m.Session, *e.Error)
	// Update some details of session, e.g. expiration time

	UpdateSession(session *m.Session) *e.Error
}

// A new simple session service that contains basic functionalities.
type sSessionService struct {
	session dal.SessionDAL
	user    dal.UserDAL
	logger  l.Logger
}

func newSSessionService(session dal.SessionDAL, user dal.UserDAL, logger l.Logger) SessionService {
	return &sSessionService{
		session,
		user,
		logger,
	}
}

type loginType int

const (
	// Just by entering phone number, we can login and create a session.
	LoginJustByPhone loginType = iota
)

func (s *sSessionService) DeleteSession(jwt *m.JWT) *e.Error {
	isMatched, err := s.session.IsMatchSessionUserID(*jwt.JTI, *jwt.UserID)
	if err != nil {
		return e.NewErrorP("failed to check if session id %s matches user id %s. (%s)",
			SEDBError, jwt.JTI, jwt.UserID, err.Error())
	}
	if !isMatched {
		return e.NewErrorP("session id %s does not belongs to user id %s", SENotFound, jwt.JTI, jwt.UserID)
	}
	if err := s.session.DeleteSession(*jwt.JTI); err != nil {
		return e.NewErrorP("failed to delete session id %s. (%s)", SEDBError, jwt.JTI, err.Error())
	}
	return nil
}

func (s *sSessionService) CreateSessionJustByPhone(session *m.Session, phone m.PhoneNumber) (*m.ID, *e.Error) {
	// Validate if the phone number is valid then return full details of the user.
	user, err := s.user.GetUserByID(session.UserID)
	if err != nil {
		return nil, e.NewErrorP("failed to get user by its id %s. (%s)", SEDBError, session.UserID, err.Error())
	} else if user == nil {
		return nil, e.NewErrorP("user with id %s not found", SENotFound, session.UserID.ToString())
	} else if user.PhoneNumber != phone {
		return nil, e.NewErrorP("user with id %s does not match the phone number", SENotFound, session.UserID.ToString())
	}
	expiredTime, _ := strconv.ParseInt(os.Getenv("SESSION_EXPIRED_TIME_MIN"), 10, 64)

	session = &m.Session{
		UserAgent:   session.UserAgent,
		UserID:      user.ID,
		IssuedAt:    time.Now().Unix(),
		ExpiredAt:   time.Now().Add(time.Minute * time.Duration(expiredTime)).Unix(),
		LastUsageAt: 0,
	}
	sessionID, err := s.session.CreateSession(session)
	if err != nil {
		return nil, e.NewErrorP("failed to create session for user id %s. (%s)", SEDBError, session.UserID, err.Error())
	}
	return sessionID, nil
}

func (s *sSessionService) ValidateSessionJWT(token string) (*m.JWT, *e.Error) {
	validJWT, validationErr := s.validateJWT(token)
	if validationErr != nil {
		return nil, validationErr.AppendBegin("failed to validate JWT token")
	}
	// Check that the session has not been deleted.
	session, err := s.session.GetSessionByID(*validJWT.JTI)
	if err != nil {
		return nil, e.NewErrorP("failed to get session id %s. (%s)", SEDBError, validJWT.JTI, err.Error())
	} else if session == nil {
		return nil, e.NewErrorP("session with id %s not found", SENotFound, validJWT.JTI)
	}
	return validJWT, nil
}

func (s *sSessionService) validateJWT(token string) (*m.JWT, *e.Error) {
	if token == "" {
		return nil, e.NewErrorP("JWT token is empty", SEAuthFailed)
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// It's type assertion
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, e.NewSError("JWT token is invalid")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEAuthFailed)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		sub, okSub := claims["sub"].(string)
		jp, okJP := claims["jp"].(string)
		iat, okIat := claims["iat"].(int64)
		exp, okExp := claims["exp"].(int64)
		jti, okJti := claims["jti"].(string)
		userID, userErr := m.ID{}.FromString(sub)
		jPID, jpErr := m.ID{}.FromString(jp)
		jtiID, jtiErr := m.ID{}.FromString(jti)

		if okSub && okIat && okExp && okJP && okJti && userErr == nil && jpErr == nil && jtiErr == nil {
			// Validate time of token such that IAT < now < EXP
			if err := validateJWTTime(iat, exp); err != nil {
				return nil, e.NewErrorP(err.Error(), SEAuthFailed)
			}
			return &m.JWT{
				UserID: &userID,
				JPID:   &jPID,
				JTI:    &jtiID,
				IAT:    iat,
				EXP:    exp,
			}, nil
		}
		return nil, e.NewErrorP("JWT token is invalid. (sub: %s, jp: %s, iat: %d, exp: %d",
			SEAuthFailed, sub, jp, iat, exp)
	} else {
		return nil, e.NewErrorP("JWT token is invalid", SEAuthFailed)
	}
}

func (s *sSessionService) UpdateSession(session *m.Session) *e.Error {
	panic("UpdateSessionExp doesn't implemented")
}

// Validate Issued At Timestamp and Expiration timestamp
func validateJWTTime(iat, exp int64) error {
	issuedAt := time.Unix(iat, 0)
	if time.Now().UTC().Before(issuedAt) {
		return fmt.Errorf("IAT token issued in the future. IAT represents %s", issuedAt.Format(time.RFC3339))
	}
	expiredAt := time.Unix(exp, 0)
	if time.Now().UTC().After(expiredAt) {
		return fmt.Errorf("EXP token expired. EXP represents %s", expiredAt.Format(time.RFC3339))
	}
	return nil
}

func (s *sSessionService) GetSessionByID(sessionID *m.ID) (*m.Session, *e.Error) {
	session, err := s.session.GetSessionByID(*sessionID)
	if err != nil {
		return nil, e.NewErrorP("failed to get session id %s. (%s)", SEDBError, sessionID, err.Error())
	}
	return session, nil
}
