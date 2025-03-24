package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type EventHttp struct {
	eventService s.EventService
	logger       l.Logger
}

func newEventHttp(event s.EventService, logger l.Logger) EventHttp {
	return EventHttp{
		eventService: event,
		logger:       logger,
	}
}

// @Security BearerAuth
// @Summary Create event
// @Description Create event for specified job position and return its id.
// @Tags event
// @Accept json
// @Produce json
// @Param event body models.Event true "Event"
// @Success 200 {object} HttpResponse{details=idResponse} "Success creating event"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 404 {object} HttpResponse{details=string} "Not found error. The job position doesn't belongs to current user."
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /events [post]
func (h *EventHttp) CreateEvent(c *gin.Context) {
	event := m.Event{}
	if err := parseValidateJSON(c, &event, h.logger); err != nil {
		return
	}
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}

	id, err := h.eventService.CreateEvent(event, jwt.UserID)
	if err == nil {
		h.logger.Debugf("Created event with id %s.", id.ToString())
		successResp(c, MsgEventCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create event (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		h.logger.Debugf("Failed to create event: %s", err.Error())
		notFoundResp(c, fmt.Sprintf(MsgNotFoundC, MsgJP), MsgCheckInfoAgain)
	default:
		h.logger.Panicf("Unexpected error code %d (%s)", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
