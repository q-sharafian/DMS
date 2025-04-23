package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
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
	h.logger.Debugf("Got session jwt token. (jwt token length: %d)", len(jwt))
	jwt = strings.Replace(jwt, "Bearer ", "", 1)
	jwt = strings.TrimSpace(jwt)

	params, err := h.sessionService.ValidateSessionJWT(m.Token(jwt))
	if err == nil {
		c.Set(authInfo, params)
		c.Next()
		return
	}
	c.Abort()
	switch code := err.GetCode(); code {
	case s.SEAuthFailed:
		h.logger.Debugf("Failed to authenticate user in the middleware: %s", err.Error())
		unauthorizedResp(c, MsgAuthFailed, MsgTryAgain)
	case s.SEDBError:
		h.logger.Errorf("Failed to authenticate user (%s)", err.Error())
		unauthorizedResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		h.logger.Debugf("Failed to authenticate user in the middleware: %s", err.Error())
		unauthorizedResp(c, MsgAuthNotFound, MsgReferAdmin)
	default:
		h.logger.Panicf("Error code %d doesn't handled. (%s)", code, err.Error())
	}
}

var corsConfig cors.Config
var isCorsRunned = false

// Cors enables Cross-Origin Resource Sharing (CORS) headers for the given
// request.
func (h MiddlewareHttp) Cors(c *gin.Context) {
	if !isCorsRunned {
		allowedOriginsEnv := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), " ")
		allowedOrigins := make([]string, 0, len(allowedOriginsEnv))
		for _, origin := range allowedOriginsEnv {
			if strings.TrimSpace(origin) == "" {
				continue
			}
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
		}
		isCorsRunned = true
		h.logger.Debugf("Allowed origins: %+v", allowedOrigins)

		corsConfig = cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * 60 * 60,
		}
	}

	corsMiddleware := cors.New(corsConfig)
	corsMiddleware(c)
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

// TODO: Create a controller that abort requests if the specified user is disabled.
