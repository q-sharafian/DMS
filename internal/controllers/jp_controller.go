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

func (h *JPHttp) CreateJP(c *gin.Context) {
	jp := struct {
		JobPosition m.JobPosotion
		Permission  m.Permission
	}{
		JobPosition: m.JobPosotion{},
		Permission:  m.Permission{},
	}
	if err := parseValidateJSON(c, &jp, h.logger); err != nil {
		return
	}

	id, err := h.jpService.CreateJP(&jp.JobPosition, &jp.Permission)
	if err == nil {
		successResp(c, MsgJPCreated, "")
		h.logger.Debugf("Created job position with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create job position (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

func (h *JPHttp) GetUserJPs(c *gin.Context) {
	paramParser := newParamParser(c, h.logger)
	var userID m.ID
	if paramParser.parseID("id", &userID) != nil {
		return
	}
	jps, err2 := h.jpService.GetUserJPs(userID)
	switch code := err2.GetCode(); code {
	case s.SEDBError:
		h.logger.Infof("Failed to get job positions for user %s (%s)", userID.ToString(), err2.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	default:
		successResp(c, MsgSuccessAction, jps)
	}
}
