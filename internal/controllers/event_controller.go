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
		h.logger.Debugf("Created event with id %s.", id.String())
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

// @Security BearerAuth
// @Summary Get last N events by job position id
// @Description Get last N events by job position id.
// @Tags event
// @Accept json
// @Produce json
// @Param jpid query string true "Job position id"
// @Param limit query int false "Limit of events to fetch. Default is 40. Max is 100. if limit be equals 0, then return all events from offset to the end."
// @Param offset query int false "Offset of events to fetch. Default is 0."
// @Success 200 {object} HttpResponse{details=[]models.Event} "Success fetching events"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 403 {object} HttpResponse{details=string} "Jon position doesn't belong to current user."
// @Router /events [get]
func (h *EventHttp) GetNLastEventsByJPID(c *gin.Context) {
	queryParser := newQueryParser(c, h.logger)
	limitDefaultValue := uint64(40)
	limit, _ := queryParser.ParseUInt("limit", &limitDefaultValue)
	maxLimit := uint64(100)
	if *limit > maxLimit {
		*limit = maxLimit
	}
	offsetDefaultValue := uint64(0)
	offset, _ := queryParser.ParseUInt("offset", &offsetDefaultValue)
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}
	jpID, err := queryParser.ParseID("jpid", nil)
	if err != nil {
		h.logger.Debugf("Failed to parse job position id: %s", err.Error())
		// customErrResp(c, hCBadValue, MsgBadValue, fmt.Sprintf(MsgIsNotValidC, MsgJP))
		return
	} else if jpID.IsNil() {
		customErrResp(c, hCBadValue, MsgBadValue, fmt.Sprintf(MsgRequiredValueC, MsgJP))
		return
	}

	events, err2 := h.eventService.GetNLastEventsByJPID(jwt.UserID, *jpID, *limit, *offset)
	if err2 == nil {
		h.logger.Debugf("Fetched %d events for job position id %s. (limit: %d, offset: %d)",
			len(*events), jpID.String(), *limit, *offset)
		successResp(c, MsgSuccessAction, events)
		return
	}

	switch code := err2.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to fetch events (%s)", err2.Error())
		customErrResp(c, hCDBError, MsgServerError, MsgTryAgain)
	case s.SEJPNotMatchedUser:
		h.logger.Debugf("Failed to fetch events: %s", err2.Error())
		customErrResp(c, hCJPNotMatchedUser, fmt.Sprintf(MsgNotFoundC, MsgJP), MsgCheckInfoAgain)
	default:
		h.logger.Panicf("Unexpected error code %d (%s)", code, err2.Error())
		customErrResp(c, hcUnexpectedError, MsgServerError, MsgTryAgain)
	}
}
