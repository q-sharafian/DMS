package controllers

import (
	e "DMS/internal/error"
	l "DMS/internal/logger"
	s "DMS/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// init function in GO runs once per package within a module before the main function.
func init() {
	validate = validator.New()
}

type HttpResponse struct {
	// Its type of response. e.g. error, warning, or success.
	// It's better to just have these three types.
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// All http controllers combine to this struct with some related functionalities.
type HttpConrtoller struct {
	User UserHttp
	// JP     JPHttp
	logger l.Logger
}

// Create new controller layer based on HTTP protocol
func NewHttpController(services s.Service, logger l.Logger) HttpConrtoller {
	return HttpConrtoller{
		User: newUserHttp(services.User, logger),
		// JP:     newJPHttp(services.JP, logger),
		logger: logger,
	}
}

// Postfix "C" means it's more complete comment compared to its pair without postfix "C"
// and these complete comment have format specifier and must be used with something
// like fmt.Printf() functions.
const (
	BadJsonStruct          = "ساختار داده ورودی اشتباه است"
	ParsingError           = "خطایی در هنگام تجزیه رخ داد"
	ServerError            = "خطایی در سمت سرور رخ داد"
	DisabledUserC          = "کاربر %s غیر فعال شده است"
	DisabledUsed           = "کاربر غیر فعال شده است"
	FixDisabledUserProblem = "جهت رفع اشکال به سرپرست مراجعه کنید"
	UserExists             = "کاربر از قبل وجود دارد"
	UserExistsExpanded     = "کاربری با این شماره تلفن که قصد ساخت آن را دارید از قبل وجود دارد"
	TryAgain               = "لطفا مجددا تلاش نمایید"
	UserCreated            = "کاربر جدید با موفقیت ساخته شد"
	AdminCreated           = "مدیر جدید با موفقیت ساخته شد"
	FailedCreatingUser     = "ساخت کاربر جدید با خطا مواجه شد"
	JPCreated              = "سمت جدید با موفقیت ساخته شد"
)

func formatResponse(c *gin.Context, httpCode int, typeResp, msg, details string) {
	c.JSON(httpCode, HttpResponse{
		Type:    typeResp,
		Code:    httpCode,
		Message: msg,
		Details: details,
	})
}

// Return a JSON response with HTTP code 400 to the client
func badRequestResp(c *gin.Context, message, details string) {
	formatResponse(c, http.StatusBadRequest, "error", message, details)
}

// Return a JSON response with HTTP code 200 to the client that represents success.
func successResp(c *gin.Context, message, details string) {
	formatResponse(c, http.StatusOK, "success", message, details)
}

// Return a JSON response with HTTP code 200 to the client that represents warning.
func warningResp(c *gin.Context, message, details string) {
	formatResponse(c, http.StatusOK, "warning", message, details)
}

// Return a JSON response with HTTP code 200 to the client that represents error.
func errResp(c *gin.Context, message, details string) {
	formatResponse(c, http.StatusOK, "error", message, details)
}

// Return a JSON response with HTTP code 500 to the client
func serverErrResp(c *gin.Context, message, details string) {
	formatResponse(c, http.StatusBadRequest, "error", message, details)
}

// Try to parse and validate input object and return error if it's not valid.
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
		badRequestResp(c, BadJsonStruct, ParsingError)
		// id := c.Param("id")
		// c.Status(400)
		return e.NewSError("couldn't parse JSON object")
	}
	if err := validate.Struct(obj); err != nil {
		logger.Debugf("Error in validating struct (%s)", err.Error())
		badRequestResp(c, BadJsonStruct, ParsingError)
		return e.NewSError("couldn't validate JSON object")
	}
	return nil
}
