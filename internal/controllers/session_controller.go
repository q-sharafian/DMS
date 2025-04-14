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

// @Security BearerAuth
// @Summary Login/Create JWT with phone number only
// @Description Login/Create JWT with phone number only.
// @Tags session
// @Param phone body models.PhoneBasedLoginInfo true "Phone number"
// @Success 200 {object} HttpResponse{details=string} "Success login and response created JWT token"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 401 {object} HttpResponse{details=string} "User not found with such phone number"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /login/phone-based [post]
func (h *SessionHttp) PhoneBasedLogin(c *gin.Context) {
	session := m.PhoneBasedLoginInfo{}
	if err := parseValidateJSON(c, &session, h.logger); err != nil {
		return
	}
	token, err := h.sessionService.CreateSessionJustByPhone(&session)
	if err == nil {
		h.logger.Debugf("Created session with user-agent %s.", session.UserAgent)
		successResp(c, MsgSuccessfulLogin, token)
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create session (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		unauthorizedResp(c, MsgAuthNotFound, MsgReferAdmin)
	case s.SEEncodingError:
		h.logger.Debugf("Failed to create session: %s", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	default:
		h.logger.Errorf("Unexpected error code %d (%s)", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Security BearerAuth
// @Summary Logout
// @Description Logout from the current session.
// @Tags session
// @Success 200 {object} HttpResponse{details=string} "Success logout"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 401 {object} HttpResponse{details=string} "Unauthorized access to resource"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /logout [post]
func (h *SessionHttp) Logout(c *gin.Context) {
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}
	err := h.sessionService.DeleteSession(jwt)
	if err == nil {
		h.logger.Debugf("Deleted session with id %s.", jwt.JTI.String())
		successResp(c, MsgSuccessfulLogout, MsgSuccessfulLogout)
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to delete session (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		h.logger.Debugf("The session with id %s was not found. (%s)", jwt.JTI.String(), err.Error())
		unauthorizedResp(c, MsgSessionNotFound, MsgSessionNotFound)
	case s.SEDeletedPreviously:
		h.logger.Debugf("The session with id %s was deleted/deactivated previously. (%s)", jwt.JTI.String(), err.Error())
		unauthorizedResp(c, MsgSessionNotFound, MsgDeletedSessionPreviously)
	default:
		h.logger.Errorf("Unexpected error code %d (%s)", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
