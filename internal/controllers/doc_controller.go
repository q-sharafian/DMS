package controllers

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"

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
func (h *DocHttp) CreateDoc(ctx *gin.Context) {
	doc := m.Doc{}
	if err := parseValidateJSON(ctx, &doc, h.logger); err != nil {
		return
	}
	// It's possible that the job position field specified in the received document is
	// manipulated and not true, so we use the job position specified in the middleware.
	jwt := getJWT(ctx, h.logger)
	if jwt == nil {
		return
	}

	id, err := h.docService.CreateDoc(&doc, jwt.UserID)
	if err == nil {
		h.logger.Debugf("Created doc with id %s successfully", id.ToString())
		successResp(ctx, MsgDocCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		forbiddenErrResp(ctx, MsgDisabledUser, MsgFixDisabledUserProblem)
	case s.SEDBError:
		serverErrResp(ctx, MsgServerError, MsgTryAgain)
	case s.SENotFound:
		h.logger.Infof("Failed to create doc (%s)", err.Error())
		notFoundResp(ctx, MsgAuthNotFound, MsgReferAdmin)
	default:
		h.logger.Panicf("Unexpected error code %d in CreateDoc controller: %s", code, err.Error())
		serverErrResp(ctx, MsgServerError, MsgTryAgain)
	}
}

func (h *DocHttp) GetNLastDocsByEventID(c *gin.Context) {
	paramParser := newParamParser(c, h.logger)
	var eventID m.ID
	if paramParser.parseID("id", &eventID) != nil {
		return
	}
	queryParser := newQueryParser(c, h.logger)
	var count *uint64
	defaultValue := uint64(10)
	count, _ = queryParser.ParseUInt("count", &defaultValue)

	authInfo := getJWT(c, h.logger)
	if authInfo == nil {
		return
	}
	docs, err := h.docService.GetNLastDocByEventID(eventID, nil, authInfo.JPID, int(*count))
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to get docs for event %s (%s)", eventID.ToString(), err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	default:
		successResp(c, MsgSuccessAction, docs)
	}
}
