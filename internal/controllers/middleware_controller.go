package controllers

import (
	l "DMS/internal/logger"
	s "DMS/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiddlewareHttp struct {
	logger         l.Logger
	sessionService s.SessionService
}

func newMiddlewareHttp(sessionService s.SessionService, logger l.Logger) MiddlewareHttp {
	return MiddlewareHttp{
		logger,
		sessionService,
	}
}

func (h MiddlewareHttp) Authentication(c *gin.Context) {
	jwt := c.GetHeader("Authorization")
	h.logger.Debugf("Got jwt: \"%s\"", jwt)
	jwt = strings.Replace(jwt, "Bearer ", "", 1)
	jwt = strings.TrimSpace(jwt)

	params, err := h.sessionService.ValidateSessionJWT(jwt)
	if err == nil {
		c.Set(authInfo, params)
		c.Next()
		return
	}
	c.Abort()
	switch code := err.GetCode(); code {
	case s.SEAuthFailed:
		unauthorizedResp(c, MsgAuthFailed, MsgTryAgain)
	case s.SEDBError:
		h.logger.Errorf("Failed to authenticate user (%s)", err.Error())
		unauthorizedResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		unauthorizedResp(c, MsgAuthNotFound, MsgReferAdmin)
	default:
		h.logger.Panicf("Error code %d doesn't handled. (%s)", code, err.Error())
	}
}

// TODO: Create a controller that abort requests if the specified user is disabled.
