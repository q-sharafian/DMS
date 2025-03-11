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

func NewJPHttp(jobPositionService s.JPService, logger l.Logger) JPHttp {
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
		successResp(c, JPCreated, "")
		h.logger.Debugf("Created job position with id %s successfully", id.ToString())
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create job position (%s)", err.Error())
		serverErrResp(c, ServerError, TryAgain)
	}
}
