package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type DocHttp struct {
	docService s.DocService
	logger     l.Logger
}

func newDocHttp(docService s.DocService, logger l.Logger) DocHttp {
	return DocHttp{
		docService,
		logger,
	}
}

// @Security BearerAuth
// @Summary Create document
// @Description Create document for specified event and current user in the current time and return its id.
// @Tags document
// @Accept json
// @Produce json
// @Param doc body models.Doc true "Doc"
// @Success 200 {object} HttpResponse{details=idResponse} "Success creating document. Returns the document id."
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 404 {object} HttpResponse{details=string} "Not found error. The job position doesn't belongs to current user."
// @Failure 403 {object} HttpResponse{details=string} "Forbidden error. The user is disabled."
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /docs [post]
func (h *DocHttp) CreateDoc(c *gin.Context) {
	doc := m.Doc{}
	if err := parseValidateJSON(c, &doc, h.logger); err != nil {
		return
	}
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}

	id, err := h.docService.CreateDoc(&doc, jwt.UserID)
	if err == nil {
		h.logger.Debugf("Created doc with id %s successfully", id.String())
		successResp(c, MsgDocCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		forbiddenErrResp(c, MsgDisabledUser, MsgFixDisabledUserProblem)
	case s.SEDBError:
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		h.logger.Infof("Failed to create doc (%s)", err.Error())
		notFoundResp(c, MsgAuthNotFound, MsgReferAdmin)
	case s.SEEventOwnerMismatched:
		h.logger.Debugf("The job position %s can't create a doc for the event %s: %s",
			doc.CreatedBy.String(), doc.EventID.String(), err.Error())
		forbiddenErrResp(c, fmt.Sprintf(MsgCreationNotAllowC, MsgDocs), MsgEventOwnerMismatchedJP)
	default:
		h.logger.Panicf("Unexpected error code %d in CreateDoc controller: %s", code, err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Security BearerAuth
// @Summary Get last documents
// @Description Get the latest documents within the event with the given ID. The documents can be retrieved only by the owner of the event and its ancestors, if they have the appropriate permissions.
// @Tags document
// @Accept json
// @Produce json
// @Param jp_id path string true "Job position id"
// @Param event_id path string true "Event id"
// @Param count query int false "Number of documents to get"
// @Success 200 {object} HttpResponse{details=[]models.Doc} "Documents"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 404 {object} HttpResponse{details=string} "Not found error. The event doesn't exists."
// @Failure 403 {object} HttpResponse{details=string} "Forbidden error. The job position doesn't have permission to access this event and their docs."
// @Failure 401 {object} HttpResponse{details=string} "The user is not authenticated"
// @Failure 400 {object} HttpResponse{details=string} "Bad request error"
// @Router /jps/{jp_id}/events/{event_id}/docs [get]
func (h *DocHttp) GetNLastDocsByEventID(c *gin.Context) {
	paramParser := newParamParser(c, h.logger)
	var eventID *m.ID
	var err error
	if eventID, err = paramParser.parseID("event_id", nil); err != nil {
		return
	}
	var jPID *m.ID
	if jPID, err = paramParser.parseID("jp_id", nil); err != nil {
		return
	}
	queryParser := newQueryParser(c, h.logger)
	var count *uint64
	defaultValue := uint64(10)
	count, _ = queryParser.ParseUInt("count", &defaultValue)
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}

	h.logger.Debugf("Getting last %d docs for event %s", *count, eventID.String())
	docs, err2 := h.docService.GetNLastDocByEventID(*eventID, jwt.UserID, nil, *jPID, int(*count))
	if err2 == nil {
		h.logger.Debugf("Got last %d docs for event %s successfully", *count, eventID.String())
		successResp(c, MsgSuccessAction, docs)
		return
	}
	switch code := err2.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to get docs for event %s (%s)", eventID.String(), err2.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	case s.SEJPNotMatchedUser:
		h.logger.Debugf("Job position doesn't belong the user: %s", err2.Error())
		forbiddenErrResp(c, MsgJPNotBelongsUser, MsgReferAdmin)
	case s.SENotAncestor:
		h.logger.Debugf("Job position with id %s is not ancestor of job position who created event with id %s: %s",
			jPID.String(), eventID.String(), err2.Error())
		forbiddenErrResp(c, MsgNotPermission, MsgNotAncestor)
	default:
		h.logger.Panicf("Unexpected error code %d in GetNLastDocsByEventID controller: %s", code, err2.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}

// @Security BearerAuth
// @Summary Get n last documents that are accessible for the user who sent the request.
// @Description Get n last documents (according to the limit and offset values) that are accessible for the user (and one of his job positions) who sent the request. (If the job position is admin, he has access to all documents.) For example, offset = 20 and limit = 10 would return 10 records (documents 21-30).
// @Tags document
// @Accept json
// @Produce json
// @Param limit query int false "Number of documents to get. Maximum is 50"
// @Param offset query int false "Number of documents to skip"
// @Param jpid query string true "Job position id"
// @Success 200 {object} HttpResponse{details=[]models.DocWithSomeDetails} "Documents"
// @Failure 500 {object} HttpResponse{details=string} "Server or database error"
// @Failure 403 {object} HttpResponse{details=string} "Forbidden error. The user is not authorized to access this resource, job position doesn't belongs to the user or etc."
// @Failure 401 {object} HttpResponse{details=string} "The user is not authorized"
// Failure 400 {object} HttpResponse{details=string} "Bad request error/parsing error"
// @Router /docs [get]
func (h *DocHttp) GetNLastDocs(c *gin.Context) {
	queryParser := newQueryParser(c, h.logger)
	limitDefaultValue := uint64(20)
	limit, _ := queryParser.ParseUInt("limit", &limitDefaultValue)
	maxLimit := uint64(50)
	if *limit > maxLimit {
		*limit = 50
	} else if *limit < 1 {
		*limit = 1
	}

	offsetDefaultValue := uint64(0)
	offset, _ := queryParser.ParseUInt("offset", &offsetDefaultValue)
	jwt := getJWT(c, h.logger)
	if jwt == nil {
		return
	}
	jpID, err := queryParser.ParseID("jpid", nil)
	if err != nil {
		h.logger.Debugf("Error in parsing jpid: %s", err.Error())
		return
	} else if jpID.IsNil() {
		customErrResp(c, hCBadValue, MsgBadValue, fmt.Sprintf(MsgRequiredValueC, MsgJP))
		return
	}

	h.logger.Debugf("Getting last %d docs (with offset %d)", *limit, *offset)
	docs, err2 := h.docService.GetNLastDocs(jwt.UserID, *jpID, *limit, *offset)
	if err2 == nil {
		h.logger.Debugf("Got last %d docs (with offset %d) successfully", *limit, *offset)
		successResp(c, MsgSuccessAction, *docs)
		return
	}
	switch code := err2.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to get last %d docs (with offset %d): %s", *limit, *offset, err2.Error())
		customErrResp(c, hCDBError, MsgServerError, MsgTryAgain)
	case s.SEJPNotMatchedUser:
		h.logger.Debugf("Job position doesn't belong the user: %s", err2.Error())
		customErrResp(c, hCJPNotMatchedUser, MsgJPNotBelongsUser, MsgReferAdmin)
	default:
		h.logger.Panicf("Unexpected error code %d in GetNLastDocs controller: %s", code, err2.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	}
}
