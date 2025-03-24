package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"

	"github.com/gin-gonic/gin"
)

// Job position controller
type JPHttp struct {
	jpService s.JPService
	logger    l.Logger
}

func newJPHttp(jobPositionService s.JPService, logger l.Logger) JPHttp {
	return JPHttp{jobPositionService, logger}
}

// @Summary Create a new user job position
// @Description Create a new job position for specified user. Each user job position must be created with another job position.
// @Tags job-position
// @Param jPWithPermission body models.UserJPWithPermission true "Job position"
// @Success 200 {object} HttpResponse{details=idResponse} "Job position created and response its id"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Failure 401 {object} HttpResponse{details=string} "Unauthorized access to resource"
// @Router /jps [post]
func (h *JPHttp) CreateUserJP(c *gin.Context) {
	jp := m.UserJPWithPermission{
		JobPosition: m.UserJobPosition{},
		Permission:  m.Permission{},
	}
	if err := parseValidateJSON(c, &jp, h.logger); err != nil {
		return
	}
	h.logger.Debugf("Got job position %+v and permission %+v, ParentID: %+v", jp.JobPosition, jp.Permission, jp.JobPosition.ParentID)
	id, err := h.jpService.CreateUserJP(&jp.JobPosition, &jp.Permission)
	if err == nil {
		successResp(c, MsgJPCreated, newIDResponse(*id))
		h.logger.Debugf("Created job position with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create job position (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Summary Create a new job position
// @Description Create a new job position for specified user. Each Admin job position is created without a job position and has no parent job position.
// @Tags job-position
// @Param jPWithPermission body models.AdminJPWithPermission true "Job position"
// @Success 200 {object} HttpResponse{details=idResponse} "Job position created and response its id"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Failure 401 {object} HttpResponse{details=string} "Unauthorized access to resource"
// @Router /jps/admin [post]
func (h *JPHttp) CreateAdminJP(c *gin.Context) {
	jp := m.AdminJPWithPermission{
		JobPosition: m.AdminJobPosition{},
		Permission:  m.Permission{},
	}
	if err := parseValidateJSON(c, &jp, h.logger); err != nil {
		return
	}
	h.logger.Debugf("Got job position %+v and permission %+v", jp.JobPosition, jp.Permission)
	id, err := h.jpService.CreateAdminJP(&jp.JobPosition, &jp.Permission)
	if err == nil {
		successResp(c, MsgJPCreated, newIDResponse(*id))
		h.logger.Debugf("Created job position with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create job position (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Summary Get user job positions
// @Description Get user job positions
// @Tags job-position
// @Param id query string false "User ID"
// @Param phone query string false "User phone number"
// @Success 200 {object} HttpResponse{details=[]models.UserJobPosition} "Job positions"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Failure 401 {object} HttpResponse{details=string} "Unauthorized access to resource"
// Failure 404 {object} HttpResponse{details=string} "Job positions not found"
// @Router /user/jps [get]
func (h *JPHttp) GetUserJPs(c *gin.Context) {
	queryParser := newQueryParser(c, h.logger)
	var userID *m.ID
	var err error
	if userID, err = queryParser.ParseID("id", &m.NilID); err != nil {
		h.logger.Debugf("Error in parsing user id (%s)", err.Error())
	}
	var phone *m.PhoneNumber
	if phone, err = queryParser.ParsePhone("phone", &m.NilPhone); err != nil {
		h.logger.Debugf("error in parsing user phone (%s)", err.Error())
	}

	user := m.User{}
	if userID != nil && !userID.IsNil() {
		user.ID = *userID
	}
	if phone != nil && !phone.IsNil() {
		user.PhoneNumber = *phone
	}
	jps, err2 := h.jpService.GetUserJPs(&user)
	if err2 == nil {
		successResp(c, MsgSuccessAction, jps)
		return
	}
	switch code := err2.GetCode(); code {
	case s.SEDBError:
		h.logger.Infof("Failed to get job positions for user %s (%s)", userID.ToString(), err2.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		h.logger.Debugf("Failed to get job positions for user %s (%s)", userID.ToString(), err2.Error())
		notFoundResp(c, MsgJPsNotFound, MsgCheckInfoAgain)
	default:
		h.logger.Panicf("Unexpected error code %d in GetUserJPs controller(%s)", code, err2.Error())
	}
}
