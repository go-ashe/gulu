package httpx

import (
	"errors"
	"fmt"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
)

type Error struct {
	Code     int               `json:"code"`     // 错误码
	Reason   string            `json:"reason"`   // 错误原因
	Message  string            `json:"message"`  // 前端展示错误消息
	Metadata map[string]string `json:"metadata"` // 元数据
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", e.Code, e.Reason, e.Message, e.Metadata)
}

func (e *Error) GetMessage() string {
	return e.Message
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	return &Error{
		Code:     e.Code,
		Message:  e.Message,
		Reason:   e.Reason,
		Metadata: md,
	}
}

// WithMessage with a message.
func (e *Error) WithMessage(msg string) *Error {
	return &Error{
		Code:     e.Code,
		Message:  msg,
		Reason:   e.Reason,
		Metadata: e.Metadata,
	}
}

// WithReason with a reason.
func (e *Error) WithReason(reason string) *Error {
	return &Error{
		Code:     e.Code,
		Message:  e.Message,
		Reason:   reason,
		Metadata: e.Metadata,
	}
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Reason:  reason,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...interface{}) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...interface{}) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Code returns the http code for a error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200
	}
	return FromError(err).Code
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if se := new(Error); errors.As(err, &se) {
		return se
	}

	return New(UnknownCode, UnknownReason, err.Error())
}
