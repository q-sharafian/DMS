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

// @Security BearerAuth
// @Summary Create user
// @Description Create a user and return its id. Each user must created by a job position.
// @Tags user
// @Accept json
// @Produce json
// @Param admin body models.User true "User"
// @Success 200 {object} HttpResponse{details=idResponse} "Success creating admin"
// @Failure 409 {object} HttpResponse{details=string} "This user exists previously or disabled"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /users/ [post]
func (h *UserHttp) CreateUser(c *gin.Context) {
	user := m.User{}
	if err := parseValidateJSON(c, &user, h.logger); err != nil {
		return
	}
	id, err := h.userService.CreateUser(user.Name, user.PhoneNumber, *user.CreatedBy)
	if err == nil {
		h.logger.Debugf("Created user with id %s successfully", id.String())
		successResp(c, MsgUserCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		conflictErrResp(c, MsgDisabledUser, MsgFixDisabledUserProblem)
	case s.SEExists:
		conflictErrResp(c, MsgUserExists, MsgUserExistsExpanded)
	case s.SEDBError:
		h.logger.Infof("Failed to create user with phone number %s (%s)", user.PhoneNumber, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Security BearerAuth
// @Summary Create admin
// @Description Create admin user and return its id. Admin users are users that don't have created by any user.
// @Tags user
// @Accept json
// @Produce json
// @Param adminUser body models.AdminUser true "AdminUser"
// @Success 200 {object} HttpResponse{details=idResponse} "Success creating admin"
// @Failure 409 {object} HttpResponse{details=string} "This user exists previously or disabled"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /users/admin [post]
func (h *UserHttp) CreateAdmin(c *gin.Context) {
	user := m.AdminUser{}
	if err := parseValidateJSON(c, &user, h.logger); err != nil {
		return
	}
	id, err := h.userService.CreateAdmin(user.Name, user.PhoneNumber)
	if err == nil {
		successResp(c, MsgAdminCreated, newIDResponse(*id))
		h.logger.Debugf("Created admin with id %s successfully", id.String())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		conflictErrResp(c, MsgDisabledUser, MsgFixDisabledUserProblem)
	case s.SEExists:
		conflictErrResp(c, MsgUserExists, MsgUserExistsExpanded)
	case s.SEDBError:
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Security BearerAuth
// @Summary Get details of current user
// @Description Get details of current user according to the authentication token.
// @Tags user
// @Success 200 {object} HttpResponse{details=models.User} "User details"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 404 {object} HttpResponse{details=string} "User not found"
// @Router /users/ [get]
func (h *UserHttp) GetCurrentUserInfo(c *gin.Context) {
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}
	user, err := h.userService.GetUserByID(jwt.UserID)
	if err == nil {
		successResp(c, MsgSuccessfulLogin, user)
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to get user info (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		unauthorizedResp(c, MsgAuthNotFound, MsgReferAdmin)
	default:
		h.logger.Panicf("Unexpected error %d: %s", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
