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

func (h *DocHttp) CreateDoc(ctx *gin.Context) {
	doc := m.Doc{}
	if err := parseValidateJSON(ctx, &doc, h.logger); err != nil {
		return
	}
	// It's possible that the job position field specified in the received document is
	// manipulated and not true, so we use the job position specified in the middleware.
	authInfo := getJWT(ctx, h.logger)
	if authInfo == nil {
		return
	}
	doc.CreatedBy = *authInfo.JPID

	id, err := h.docService.CreateDoc(&doc)
	if err == nil {
		h.logger.Debugf("Created doc with id %s successfully", id.ToString())
		successResp(ctx, MsgDocCreated, newIDResponse(*id))
		return
	}
	switch code := err.GetCode(); code {
	case s.SEIsDisabled:
		errResp(ctx, MsgDisabledUsed, MsgFixDisabledUserProblem)
	case s.SEDBError:
		serverErrResp(ctx, MsgServerError, MsgTryAgain)
	default:
		h.logger.Errorf("Unexpected error code %d (%s)", code, err.Error())
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
	var count uint64
	if queryParser.parseUInt("count", &count) != nil {
		return
	}

	authInfo := getJWT(c, h.logger)
	if authInfo == nil {
		return
	}
	docs, err := h.docService.GetNLastDocByEventID(eventID, nil, *authInfo.JPID, int(count))
	switch code := err.GetCode(); code {
	case s.SEDBError:
		h.logger.Errorf("Failed to get docs for event %s (%s)", eventID.ToString(), err.Error())
		serverErrResp(c, MsgServerError, MsgTryAgain)
	default:
		successResp(c, MsgSuccessAction, docs)
	}
}
