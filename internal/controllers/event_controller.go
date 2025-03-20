package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"

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

func (h *EventHttp) CreateEvent(c *gin.Context) {
	event := m.Event{}
	if err := parseValidateJSON(c, &event, h.logger); err != nil {
		return
	}
	id, err := h.eventService.CreateEvent(event)
	if err == nil {
		h.logger.Debugf("Created event with id %s.", id.ToString())
		successResp(c, MsgEventCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to create event (%s)", err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	default:
		h.logger.Errorf("Unexpected error code %d (%s)", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
