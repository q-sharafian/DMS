package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"

	"github.com/gin-gonic/gin"
)

type SessionHttp struct {
	sessionService s.SessionService
	logger         l.Logger
}

func newSessionHttp(sessionService s.SessionService, logger l.Logger) SessionHttp {
	return SessionHttp{
		sessionService: sessionService,
		logger:         logger,
	}
}

// In this login method, we login just with phone number.
func (h *SessionHttp) PhoneBasedLogin(c *gin.Context) {
	type phoneBasedLogin struct {
		m.Session
		PhoneNumber m.PhoneNumber `json:"phone_number" validate:"required"`
	}
	session := phoneBasedLogin{}
	if err := parseValidateJSON(c, &session, h.logger); err != nil {
		return
	}

	id, err := h.sessionService.CreateSessionJustByPhone(&session.Session, session.PhoneNumber)
	if err == nil {
		h.logger.Debugf("Created session with id %s.", id.ToString())
		successResp(c, MsgSuccessfulLogin, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create session (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		unauthorizedResp(c, MsgAuthNotFound, MsgReferAdmin)
	default:
		h.logger.Errorf("Unexpected error code %d (%s)", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

func (h *SessionHttp) Logout(c *gin.Context) {
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}
	err := h.sessionService.DeleteSession(jwt)
	if err == nil {
		h.logger.Debugf("Deleted session with id %s.", jwt.JTI.ToString())
		successResp(c, MsgSuccessfulLogout, "")
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to delete session (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	default:
		h.logger.Errorf("Unexpected error code %d (%s)", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
