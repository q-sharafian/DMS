package controllers

import (
	s "DMS/internal/services"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// init function in GO runs once per package within a module before the main function.
func init() {
	validate = validator.New()
}

type HttpErrResponse struct {
	ErrCode int    `json:"err_code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// All http controllers combine to this struct
type HttpConrtoller struct {
	User UserHttp
}

// Create new controller layer based on HTTP protocol
func NewHttpController(services s.Service) HttpConrtoller {
	return HttpConrtoller{
		User: newUserHttp(services.User),
	}
}

const (
	BadJsonStruct = "ساختار داده ورودی اشتباه است."
	ParsingError  = "خطایی در هنگام تجزیه رخ داد"
)
