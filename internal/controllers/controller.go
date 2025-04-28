package controllers

import (
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	s "DMS/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// init function in GO runs once per package within a module before the main function.
func init() {
	validate = validator.New()
	validate.RegisterValidation("uuidv4", m.ValidateUUIDv4)
}

type HttpResponse struct {
	// Its type of response. e.g. error, warning, or success.
	// It's better to just have these three types.
	Type string `json:"type" enums:"error,warning,success"`
	// HTTP status code
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
// like fmt.Sprintf() function.
const (
	MsgBadJsonStruct            = "ساختار داده ورودی اشتباه است"
	MsgBadValue                 = "مقدار ورودی اشتباه است"
	MsgParsingErrorC            = "خطایی در هنگام تجزیه رخ داد: %s"
	MsgParsingError             = "خطایی در هنگام تجزیه رخ داد"
	MsgServerError              = "خطایی در سمت سرور رخ داد"
	MsgDisabledUserC            = "کاربر %s غیر فعال شده است"
	MsgDisabledUser             = "کاربر غیر فعال شده است"
	MsgFixDisabledUserProblem   = "جهت رفع اشکال به سرپرست مراجعه کنید"
	MsgUserExists               = "کاربر از قبل وجود دارد"
	MsgUserExistsExpanded       = "کاربری با این شماره تلفن که قصد ساخت آن را دارید از قبل وجود دارد"
	MsgTryAgain                 = "لطفا مجددا تلاش نمایید"
	MsgUserCreated              = "کاربر جدید با موفقیت ساخته شد"
	MsgAdminCreated             = "مدیر جدید با موفقیت ساخته شد"
	MsgFailedCreatingUser       = "ساخت کاربر جدید با خطا مواجه شد"
	MsgJPCreated                = "سمت جدید با موفقیت ساخته شد"
	MsgSuccessAction            = "عملیات با موفقیت انجام شد"
	MsgEventCreated             = "رویداد با موفقیت ایجاد شد"
	MsgDocCreated               = "مستند با موفقیت ایجاد شد"
	MsgAuthFailed               = "احراز هویت ناموفق"
	MsgAuthNotFound             = "کاربر و یا سمت شغلی مورد نظر موجود نیست"
	MsgReferAdmin               = "برای رفع اشکال پشتیبانی مراجعه کنید"
	MsgSuccessfulLogin          = "ورود با موفقیت انجام شد"
	MsgSuccessfulLogout         = "با موفقیت از حساب خارج شدید"
	MsgSessionNotFound          = "جلسه مورد نظر یافت نشد"
	MsgDeletedSessionPreviously = "جلسه از قبل غیر فعال شده است"
	MsgJPsNotFound              = "سمت شغلی یافت نشد"
	MsgCheckInfoAgain           = "لطفا مشخصات را مجددا بررسی کنید"
	MsgNotFoundC                = "مورد مورد نظر یافت نشد %s"
	MsgJP                       = "سمت شغلی"
	MsgDocs                     = "مستندات"
	MsgCreationNotAllowC        = "ایجاد %s برای شما مجاز نیست"
	MsgEventOwnerMismatchedJP   = "کسی که این رویداد را ایجاد کرده است، از این سمت شغلی استفاده نمی کند"
	MsgJPNotBelongsUser         = "عنوان شغلی مورد نظر به این کاربر تعلق ندارد"
	MsgNotAncestor              = "عنوان شغلی جاری، پایین تر از عوان شغلی مورد نظر است"
	MsgNotPermission            = "مجوز دسترسی ندارید"
	MsgSomeActionsFailed        = "خطایی در برخی بخش‌ها رخ داده است"
)

// hC = http code
const (
	hCDBError          = http.StatusInternalServerError
	hCJPNotMatchedUser = http.StatusForbidden
	hCBadValue         = http.StatusBadRequest
	hCParsingError     = http.StatusBadRequest
)
const authInfo = "AuthInfo"

// It's used to response just an id to the client
type idResponse struct {
	// Because the UUID in the response will be an array, we use string as id.
	ID string `json:"id,omitempty" example:"8b2d1c6b-6c2c-4a8b-8b2d-1c6b6c2c4a8b"`
}

func newIDResponse(id m.ID) idResponse {
	return idResponse{
		ID: id.String(),
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

// Return a JSON response with HTTP code 404 to the client that represents success.
func notFoundResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusNotFound, "error", message, details)
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

// Return a JSON response with HTTP code 403 to the client
func forbiddenErrResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusForbidden, "error", message, details)
}

// Return a JSON response with HTTP code 409 to the client
func conflictErrResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusConflict, "error", message, details)
}

// Return a JSON response with HTTP code 401 to the client
func unauthorizedResp(c *gin.Context, message string, details any) {
	formatResponse(c, http.StatusUnauthorized, "error", message, details)
}

func customErrResp(c *gin.Context, httpCode int, message string, details any) {
	formatResponse(c, httpCode, "error", message, details)
}

// Try to parse and validate input object with V10 and return error if it's not valid.
// If occured error during parsing or validating, return HTTP bad request error (code 400)
// response and create log. So we don't need to send HTTP response for error.
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

// Parse the input HTTP parameter specified with the input key (contained in url) to an ID type.
// If an error occurs during parsing or validation, create a log. If the default value
// isn't null, return the default value, but if it is null, return an error as well.
func (p *parseParam) parseID(paramKey string, defaultValue *m.ID) (*m.ID, error) {
	param := p.c.Param(paramKey)
	id, err := m.ID{}.FromString2(param)
	if err != nil {
		p.logger.Debugf("Error in parsing id \"%s\" (%s)", param, err.Error())
		if defaultValue == nil {
			badRequestResp(p.c, MsgBadValue, MsgParsingError)
		}
		return defaultValue, e.NewSError("couldn't parse ID")
	}
	return &id, nil
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

// A Parser to parse the input HTTP query specified with the input key (contains in url) to a suitable type and return it.
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

// ParseID parses the input HTTP query parameter specified by queryKey into an ID type.
// If the parameter is missing or cannot be parsed, it returns the defaultValue if provided.
// Otherwise, it logs the parsing error and returns an error. If the input is invalid
// and no default value is provided, it sends an HTTP bad request response.
func (p *queryParser) ParseID(queryKey string, defaultValue *m.ID) (*m.ID, error) {
	param, ok := p.c.GetQuery(queryKey)
	if !ok {
		if defaultValue != nil {
			return defaultValue, nil
		}
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return nil, fmt.Errorf("the input parameter and default values are empty and nil")
	}

	id, err := m.ID{}.FromString2(param)
	if err != nil {
		p.logger.Debugf("Error in parsing id \"%s\" (%s)", param, err.Error())
		if defaultValue == nil {
			badRequestResp(p.c, MsgBadValue, MsgParsingError)
		}
		return defaultValue, e.NewSError("couldn't parse ID")
	}
	return &id, nil
}

// If queryKey's value is empty or raises an error during parsing, set dest with the
// default value if it's not nil and then return error and no HTTP bad request
// response will be sent. But if it's nil, return error and HTTP bad request.
func (p *queryParser) ParseUInt(queryKey string, defaultValue *uint64) (dest *uint64, err error) {
	param, ok := p.c.GetQuery(queryKey)
	if !ok {
		if defaultValue != nil {
			return defaultValue, nil
		} else {
			badRequestResp(p.c, MsgBadValue, MsgParsingError)
			return nil, fmt.Errorf("the input parameter and default values are empty and nil")
		}
	}

	uint, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		p.logger.Debugf("the input param \"%s\" must be uint but it's not. (%s)", param, err)
		if defaultValue == nil {
			badRequestResp(p.c, MsgBadValue, MsgParsingError)
		}
		return defaultValue, fmt.Errorf("the input parameter must be uint but it's not")
	}
	return &uint, nil
}

// @Refer to ParseUInt for more details
func (p *queryParser) ParseInt(queryKey string, defaultValue *int64) (*int64, error) {
	param, ok := p.c.GetQuery(queryKey)
	if !ok {
		if defaultValue != nil {
			return defaultValue, nil
		}
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return nil, fmt.Errorf("the input parameter and default values are empty and nil.")
	}

	int, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		p.logger.Debugf("the input param \"%s\" must be int but it's not. (%s)", param, err)
		if defaultValue == nil {
			badRequestResp(p.c, MsgBadValue, MsgParsingError)
		}
		return defaultValue, e.NewSError("the input parameter must be int but it's not")
	}
	return &int, nil
}

// @Refer to ParseUInt for more details
func (p *queryParser) ParsePhone(queryKey string, defaultValue *m.PhoneNumber) (*m.PhoneNumber, error) {
	param, ok := p.c.GetQuery(queryKey)
	if !ok {
		if defaultValue != nil {
			return defaultValue, nil
		}
		badRequestResp(p.c, MsgBadValue, MsgParsingError)
		return nil, fmt.Errorf("the input parameter and default values are empty and nil")
	}

	phone := m.PhoneNumber(param)
	return &phone, nil
}

// Get parsed JWT from the authentication middleware. If there's not JWT,
// the function, sends unauthorized response (code 401) to the client and returns nil.
func getJWT(c *gin.Context, logger l.Logger) *m.JWT {
	value, exists := c.Get(authInfo)
	if !exists {
		unauthorizedResp(c, MsgAuthFailed, MsgTryAgain)
		logger.Debugf("JWT is not found in the context (getAuthInfo)")
		return nil
	}
	authInfo, ok := value.(*m.JWT)
	if !ok {
		logger.Panicf("JWT type is invalid")
	}
	return authInfo
}
