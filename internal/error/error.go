package error

import "fmt"

const SeparatorMsg = ":"

type ErrorCode any

type Error struct {
	message string
	code    ErrorCode
}

// Return all error message
func (e Error) Error() string {
	return fmt.Sprintf(e.message)
}

func (e *Error) AppendEnd(msg string) {
	e.message = fmt.Sprintf("%s%s %s", e.message, SeparatorMsg, msg)
}
func (e *Error) SetCode(code ErrorCode) {
	e.code = code
}

func (e *Error) GetCode() ErrorCode {
	return e.code
}

func NewError(msg string, code ErrorCode) Error {
	return Error{
		message: msg,
		code:    code,
	}
}
