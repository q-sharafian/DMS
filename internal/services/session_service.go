package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Contains interface for all functionalities related to sessions.
type SessionService interface {
	// Create a login for the specified user and return corresponding jwt.
	//
	// Possible error codes:
	// SEDBError- SENotFound- SEEncodingError
	CreateSessionJustByPhone(details *m.PhoneBasedLoginInfo) (*string, *e.Error)
	// Delete the session associated with the JWT. Note that the user id that sends the session deletion
	// The request must match the user id that the session is created for.
	//
	// Possible error codes:
	// SEDBError- SENotFound- SEDeletedPreviously
	DeleteSession(jwt *m.JWT) *e.Error
	// Validate session based on the input jwt token. We must remove any prefix like "Bearer " from the
	// input JWT token before calling the method wih that value.
	//
	// Possible error codes:
	// SEAuthFailed- SENotFound- SEDBError
	ValidateSessionJWT(token m.Token) (*m.JWT, *e.Error)
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
	session       dal.SessionDAL
	user          dal.UserDAL
	logger        l.Logger
	rsaPrivateKey rsa.PrivateKey
	rsaPublicKey  rsa.PublicKey
}

func newSSessionService(session dal.SessionDAL, user dal.UserDAL, logger l.Logger) SessionService {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_PRIVATE_KEY")))
	if err != nil {
		logger.Panicf("Failed to parse jwt rsa private key. (%s)", err.Error())
	}
	publicKey, err2 := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("JWT_PUBLIC_KEY")))
	if err2 != nil {
		logger.Panicf("Failed to parse jwt rsa public key. (%s)", err.Error())
	}

	return &sSessionService{
		session,
		user,
		logger,
		*privateKey,
		*publicKey,
	}
}

type loginType int

const (
	// Just by entering phone number, we can login and create a session.
	LoginJustByPhone loginType = iota
)

func (s *sSessionService) DeleteSession(jwt *m.JWT) *e.Error {
	session, err := s.session.GetSessionByID(jwt.JTI)
	if err != nil {
		return e.NewErrorP("failed to get session id %s. (%s)", SEDBError, jwt.JTI.String(), err.Error())
	} else if session == nil {
		return e.NewErrorP("session with id %s id deleted/deactivated previously", SEDeletedPreviously, jwt.JTI.String())
	}
	// Check if the current user is the owner of the session.
	if session.UserID != jwt.UserID {
		return e.NewErrorP("session id %s does not belongs to user id %s", SENotFound,
			jwt.JTI.String(), jwt.UserID.String())
	}

	result, err := s.session.DeleteSession(jwt.JTI)
	if result && err == nil {
		return nil
	} else if !result && err != nil {
		return e.NewErrorP("failed to delete session id %s. (%s)", SEDBError, jwt.JTI, err.Error())
	} else if !result && err == nil {
		return e.NewErrorP("session id %s not found", SENotFound, jwt.JTI.String())
	}
	s.logger.Panicf("Unexpected DeleteSession result. result: %+v, err: %s", result, err.Error())
	return nil
}

func (s *sSessionService) CreateSessionJustByPhone(details *m.PhoneBasedLoginInfo) (*string, *e.Error) {
	user, err := s.user.GetUserByPhone(details.PhoneNumber)
	if err != nil {
		return nil, e.NewErrorP("failed to get user by its phone %s. (%s)", SEDBError, details.PhoneNumber.ToString(), err.Error())
	} else if user == nil {
		return nil, e.NewErrorP("user with phone %s not found", SENotFound, details.PhoneNumber.ToString())
	}
	expiredTime, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRED_TIME_MIN"), 10, 64)

	session := &m.Session{
		UserAgent:   details.UserAgent,
		UserID:      user.ID,
		IssuedAt:    time.Now().Unix(),
		ExpiredAt:   time.Now().Add(time.Minute * time.Duration(expiredTime)).Unix(),
		LastUsageAt: 0,
	}
	sessionID, err := s.session.CreateSession(session)
	if err != nil {
		return nil, e.NewErrorP("failed to create session for user id %s. (%s)", SEDBError, session.UserID, err.Error())
	}
	jwt := &m.JWT{
		JTI:    *sessionID,
		UserID: user.ID,
		JPID:   m.NilID,
		IAT:    session.IssuedAt,
		EXP:    session.ExpiredAt,
	}
	if jwtStr, err := s.generateJWT(jwt); err != nil {
		return nil, err.AppendBegin("failed to encodeing JWT")
	} else {
		return &jwtStr, nil
	}
}

func (s *sSessionService) ValidateSessionJWT(token m.Token) (*m.JWT, *e.Error) {
	validJWT, validationErr := s.validateJWT(token)
	if validationErr != nil {
		return nil, validationErr.AppendBegin("failed to validate JWT token")
	}
	// Check that the session has not been deleted.
	session, err := s.session.GetSessionByID(validJWT.JTI)
	if err != nil {
		return nil, e.NewErrorP("failed to get session id %s. (%s)", SEDBError, validJWT.JTI, err.Error())
	} else if session == nil {
		return nil, e.NewErrorP("session with id %s not found", SENotFound, validJWT.JTI)
	}
	return validJWT, nil
}

func (s *sSessionService) validateJWT(token m.Token) (*m.JWT, *e.Error) {
	if token == "" {
		return nil, e.NewErrorP("JWT token is empty", SEAuthFailed)
	}
	// This method checks expiration and issuer time too. So we don't need to check both.
	parsedToken, err := jwt.Parse(token.String(), func(t *jwt.Token) (any, error) {
		// It's type assertion
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("JWT token is invalid")
		}
		return &s.rsaPublicKey, nil
	})
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEAuthFailed)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		sub, okSub := claims["sub"].(string)
		jp, okJP := claims["jp_id"].(string)
		iat, okIat := claims["iat"].(float64)
		exp, okExp := claims["exp"].(float64)
		jti, okJti := claims["jti"].(string)
		userID, userErr := m.ID{}.FromString2(sub)
		jPID, jpErr := m.ID{}.FromString2(jp)
		jtiID, jtiErr := m.ID{}.FromString2(jti)

		s.logger.Debugf("Received jwt: {sub: %s, jp_id: %s, iat: %f, exp: %f, jti: %s}. Parsed IDs: {sub:%s, jp_id: %s, jti: %s}",
			sub, jp, iat, exp, jti, userID.StringP(), jPID.StringP(), jtiID.StringP())

		if okSub && okIat && okExp && okJP && okJti && userErr == nil && jpErr == nil && jtiErr == nil {
			// Validate time of token such that IAT < now < EXP
			// if err := validateJWTTime(iat, exp); err != nil {
			// 	return nil, e.NewErrorP(err.Error(), SEAuthFailed)
			// }
			return &m.JWT{
				UserID: userID,
				JPID:   jPID,
				JTI:    jtiID,
				IAT:    int64(iat),
				EXP:    int64(exp),
			}, nil
		}
		return nil, e.NewErrorP("JWT token is invalid. jwt contents=> {sub: %s, jp: %s, iat: %d, exp: %d}", SEAuthFailed, sub, jp, iat, exp)
	} else {
		return nil, e.NewErrorP("JWT token is invalid", SEAuthFailed)
	}
}

func (s *sSessionService) UpdateSession(session *m.Session) *e.Error {
	panic("UpdateSessionExp doesn't implemented")
}

func (s *sSessionService) GetSessionByID(sessionID *m.ID) (*m.Session, *e.Error) {
	session, err := s.session.GetSessionByID(*sessionID)
	if err != nil {
		return nil, e.NewErrorP("failed to get session id %s. (%s)", SEDBError, sessionID, err.Error())
	}
	return session, nil
}

// Possible error codes:
// SEEncodingError
func (s *sSessionService) generateJWT(j *m.JWT) (string, *e.Error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":     j.UserID.String(),
		"jp_id":   j.JPID.StringP(),
		"iat":     j.IAT,
		"exp":     j.EXP,
		"jti":     j.JTI.String(),
		"user_id": j.UserID.String(),
	})

	tokenString, err := token.SignedString(&s.rsaPrivateKey)
	if err != nil {
		return "", e.NewErrorP("failed to generate JWT token. (%s)", SEEncodingError, err.Error())
	}
	return tokenString, nil
}
