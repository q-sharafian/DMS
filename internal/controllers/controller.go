package controllers

import (
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

// init function in GO runs once per package within a module before the main function.
func init() {
	validate = validator.New()
	validate.RegisterValidation("uuidv4", validateUUIDv4)
}

func validateUUIDv4(fl validator.FieldLevel) bool {
	if _, err := uuid.Parse(fl.Field().String()); err != nil {
		return false
	}
	return true
}

type HttpResponse struct {
	// Its type of response. e.g. error, warning, or success.
	// It's better to just have these three types.
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// All http controllers combine to this struct with some related functionalities.
type HttpConrtoller struct {
	User       UserHttp
	JP         JPHttp
	Event      EventHttp
	Doc        DocHttp
	Middleware MiddlewareHttp
	Session    SessionHttp
	logger     l.Logger
}

// Create new controller layer based on HTTP protocol
func NewHttpController(services s.Service, logger l.Logger) HttpConrtoller {
	return HttpConrtoller{
		User:       newUserHttp(services.User, logger),
		JP:         newJPHttp(services.JP, logger),
		Event:      newEventHttp(services.Event, logger),
		Doc:        newDocHttp(services.Doc, logger),
		Middleware: newMiddlewareHttp(services.Session, logger),
		Session:    newSessionHttp(services.Session, logger),
		logger:     logger,
	}
}

// Postfix "C" means it's more complete comment compared to its pair without postfix "C"
// and these complete comment have format specifier and must be used with something
// like fmt.Printf() functions.
const (
	MsgBadJsonStruct          = "ساختار داده ورودی اشتباه است"
	MsgBadValue               = "مقدار ورودی اشتباه است"
	MsgParsingError           = "خطایی در هنگام تجزیه رخ داد"
	MsgServerError            = "خطایی در سمت سرور رخ داد"
	MsgDisabledUserC          = "کاربر %s غیر فعال شده است"
	MsgDisabledUsed           = "کاربر غیر فعال شده است"
	MsgFixDisabledUserProblem = "جهت رفع اشکال به سرپرست مراجعه کنید"
	MsgUserExists             = "کاربر از قبل وجود دارد"
	MsgUserExistsExpanded     = "کاربری با این شماره تلفن که قصد ساخت آن را دارید از قبل وجود دارد"
	MsgTryAgain               = "لطفا مجددا تلاش نمایید"
	MsgUserCreated            = "کاربر جدید با موفقیت ساخته شد"
	MsgAdminCreated           = "مدیر جدید با موفقیت ساخته شد"
	MsgFailedCreatingUser     = "ساخت کاربر جدید با خطا مواجه شد"
	MsgJPCreated              = "سمت جدید با موفقیت ساخته شد"
	MsgSuccessAction          = "عملیات با موفقیت انجام شد"
	MsgEventCreated           = "رویداد با موفقیت ایجاد شد"
	MsgDocCreated             = "مستند با موفقیت ایجاد شد"
	MsgAuthFailed             = "احراز هویت ناموفق"
	MsgAuthNotFound           = "کاربر و یا سمت شغلی مورد نظر موجود نیست"
	MsgReferAdmin             = "برای رفع اشکال پشتیبانی مراجعه کنید"
	MsgSuccessfulLogin        = "ورود با موفقیت انجام شد"
	MsgSuccessfulLogout       = "با موفقیت از حساب خارج شدید"
)

const authInfo = "AuthInfo"

// It's used to response just an id to the client
type idResponse struct {
	// Because the UUID in the response will be an array, we use string as id.
	ID string `json:"id,omitempty"`
}

func newIDResponse(id m.ID) idResponse {
	return idResponse{
		ID: id.ToString(),
	}
}
func formatResponse(c *gin.Context, httpCode int, typeResp, msg string, details any) {
	c.JSON(httpCode, HttpResponse{
		Type:    typeResp,
		Code:    httpCode,
		Message: msg,
		Details: details,
	})
}

// Return a JSON response with HTTP code 400 to the client
func badRequestResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusBadRequest, "error", message, details)
}

// Return a JSON response with HTTP code 200 to the client that represents success.
func successResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusOK, "success", message, details)
}

// Return a JSON response with HTTP code 200 to the client that represents warning.
func warningResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusOK, "warning", message, details)
}

// Return a JSON response with HTTP code 200 to the client that represents error.
func errResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusOK, "error", message, details)
}

// Return a JSON response with HTTP code 500 to the client
func serverErrResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusBadRequest, "error", message, details)
}

// Return a JSON response with HTTP code 401 to the client
func unauthorizedResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusUnauthorized, "error", message, details)
}

// Try to parse and validate input object and return error if it's not valid.
// If occured error during parsing or validating, return HTTP error response and create log.
// So we don't need to send HTTP response for error.
// Example of how to use it:
//
//	// User is a struct with some fields. e.g. ID with type int
//	user := User{}
//
//	if err := parseValidateJSON(c, &user, logger); err != nil {
//			// Handle error here.
//	}
//	fmt.Println(user) // Now, user contains the details received from client and we can use that.
func parseValidateJSON(c *gin.Context, obj any, logger l.Logger) error {
	if err := c.BindJSON(obj); err != nil {
		logger.Debugf("Error in parsing JSON object (%s)", err.Error())
		badRequestResp(c, MsgBadJsonStruct, MsgParsingError)
		// id := c.Param("id")
		// c.Status(400)
		return e.NewSError("couldn't parse JSON object")
	}
	if err := validate.Struct(obj); err != nil {
		logger.Debugf("Error in validating struct (%s)", err.Error())
		badRequestResp(c, MsgBadJsonStruct, MsgParsingError)
		return e.NewSError("couldn't validate JSON object")
	}
	return nil
}

// Parse the input HTTP parameter specified with the input key (contains in url) to a suitable type and return it.
// The value of HTTP parameter key must be not empty. We suppose they're required not optional.
// If occured error during parsing or validating, return HTTP error response and create log.
// So we don't need to send HTTP response for errors.
type parseParam struct {
	c      *gin.Context
	logger l.Logger
}

func newParamParser(c *gin.Context, logger l.Logger) *parseParam {
	return &parseParam{
		c,
		logger,
	}
}

func (p *parseParam) parseID(paramKey string, dest *m.ID) error {
	param := p.c.Param(paramKey)
	if param == "" {
		p.logger.Debugf("The param %s is empty", paramKey)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must not be empty")
	}
	uuid, err := uuid.Parse(param)
	if err != nil {
		p.logger.Debugf("Error in parsing id \"%s\" (%s)", param, err.Error())
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("couldn't parse ID")
	}
	result := m.ID{}.FromUUID(uuid)
	dest = &result
	return nil
}

func (p *parseParam) parseUInt(paramKey string, dest *uint64) error {
	param := p.c.Param(paramKey)
	if param == "" {
		p.logger.Debugf("The param %s is empty", paramKey)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must not be empty")
	}
	uint, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		p.logger.Debugf("the input param \"%s\" must be uint but it's not. (%s)", param, err)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must be uint but it's not")
	}
	dest = &uint
	return nil
}

func (p *parseParam) parseInt(paramKey string, dest *int64) error {
	param := p.c.Param(paramKey)
	if param == "" {
		p.logger.Debugf("The param %s is empty", paramKey)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must not be empty")
	}
	int, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		p.logger.Debugf("the input param \"%s\" must be int but it's not. (%s)", param, err)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must be int but it's not")
	}
	dest = &int
	return nil
}

// Parse the input HTTP query specified with the input key (contains in url) to a suitable type and return it.
// If the value of HTTP query key be empty, return nil value and nil error. It means we suppose the queries are optional.
// If occured error during parsing or validating, return HTTP error response and create log.
// So we don't need to send HTTP response for errors.
type queryParser struct {
	c      *gin.Context
	logger l.Logger
}

func newQueryParser(c *gin.Context, logger l.Logger) *queryParser {
	return &queryParser{
		c,
		logger,
	}
}

func (p *queryParser) parseID(queryKey string, dest *m.ID) error {
	param := p.c.Param(queryKey)
	if param == "" {
		dest = nil
		return nil
	}
	uuid, err := uuid.Parse(param)
	if err != nil {
		p.logger.Debugf("Error in parsing id \"%s\" (%s)", param, err.Error())
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("couldn't parse ID")
	}
	result := m.ID{}.FromUUID(uuid)
	dest = &result
	return nil
}

func (p *queryParser) parseUInt(queryKey string, dest *uint64) error {
	param := p.c.Param(queryKey)
	if param == "" {
		dest = nil
		return nil
	}
	uint, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		p.logger.Debugf("the input param \"%s\" must be uint but it's not. (%s)", param, err)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must be uint but it's not")
	}
	dest = &uint
	return nil
}

func (p *queryParser) parseInt(queryKey string, dest *int64) error {
	param := p.c.Param(queryKey)
	if param == "" {
		dest = nil
		return nil
	}
	int, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		p.logger.Debugf("the input param \"%s\" must be int but it's not. (%s)", param, err)
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return e.NewSError("the input parameter must be int but it's not")
	}
	dest = &int
	return nil
}

// Get parsed JWT from the authentication middleware. If there's not JWT,
// the function, sends unauthorized response to the client and returns nil.
func getJWT(c *gin.Context, logger l.Logger) *m.JWT {
	value, exists := c.Get(authInfo)
	if !exists {
		unauthorizedResp(c, MsgAuthFailed, MsgTryAgain)
		logger.Debugf("Authentication information is not found in the context (getAuthInfo)")
		return nil
	}
	authInfo, ok := value.(m.JWT)
	if !ok {
		logger.Panicf("Authentication information type is invalid")
	}
	return &authInfo
}
