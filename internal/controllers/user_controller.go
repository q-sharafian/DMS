package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHttp struct {
	userService s.UserService
	logger      l.Logger
}

func newUserHttp(userService s.UserService, logger l.Logger) UserHttp {
	return UserHttp{userService, logger}
}

func (h *UserHttp) CreateUser(c *gin.Context) {
	user := m.User{}
	if err := parseValidateJSON(c, &user, h.logger); err != nil {
		return
	}
	id, err := h.userService.CreateUser(user.Name, user.PhoneNumber, user.CreatedBy)
	if err == nil {
		successResp(c, UserCreated, "")
		h.logger.Debugf("Created user with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		errResp(c, DisabledUsed, FixDisabledUserProblem)
	case s.SEExists:
		errResp(c, UserExists, UserExistsExpanded)
	case s.SEDBError:
		serverErrResp(c, ServerError, TryAgain)
	}
}

func (h *UserHttp) CreateAdmin(c *gin.Context) {
	user := m.User{}
	if err := parseValidateJSON(c, &user, h.logger); err != nil {
		return
	}
	id, err := h.userService.CreateAdmin(user.Name, user.PhoneNumber)
	if err == nil {
		successResp(c, AdminCreated, "")
		h.logger.Debugf("Created admin with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		errResp(c, DisabledUsed, FixDisabledUserProblem)
	case s.SEExists:
		errResp(c, UserExists, UserExistsExpanded)
	case s.SEDBError:
		serverErrResp(c, ServerError, TryAgain)
	}
}
