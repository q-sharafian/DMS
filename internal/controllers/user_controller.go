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
		h.logger.Debugf("Created user with id %s successfully", id.ToString())
		successResp(c, MsgUserCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		errResp(c, MsgDisabledUsed, MsgFixDisabledUserProblem)
	case s.SEExists:
		errResp(c, MsgUserExists, MsgUserExistsExpanded)
	case s.SEDBError:
		h.logger.Infof("Failed to create user with phone number %s (%s)", user.PhoneNumber, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

func (h *UserHttp) CreateAdmin(c *gin.Context) {
	user := m.User{}
	if err := parseValidateJSON(c, &user, h.logger); err != nil {
		return
	}
	id, err := h.userService.CreateAdmin(user.Name, user.PhoneNumber)
	if err == nil {
		successResp(c, MsgAdminCreated, newIDResponse(*id))
		h.logger.Debugf("Created admin with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		errResp(c, MsgDisabledUsed, MsgFixDisabledUserProblem)
	case s.SEExists:
		errResp(c, MsgUserExists, MsgUserExistsExpanded)
	case s.SEDBError:
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
