package goerr

import (
	"fmt"
	"net/http"
)

type IError interface {
	Error() string
	GetCode() int
	GetDetails() []IError
	PushDetail(IError)
	GetMessage() string
	Http(code int) IError
}

type appError struct {
	code    int
	Message string   `json:"message"`
	Detail  []IError `json:"detail,omitempty"`
}

func (e *appError) PushDetail(ae IError) {
	e.Detail = append(e.Detail, ae)
}

func (e *appError) Error() (er string) {
	er += fmt.Sprintf("Code: %v; ", e.code)
	er += "Msg: " + e.Message + ";  "
	if len(e.GetDetails()) == 0 {
		return
	}
	er += " Details: {"
	for idx := range e.GetDetails() {
		er += e.GetDetails()[idx].Error()
	}
	er += "}"
	return
}

func (e *appError) GetCode() int {
	return e.code
}

func (e *appError) GetMessage() string {
	return e.Message
}

func (e *appError) GetDetails() []IError {
	return e.Detail
}

func (e *appError) Http(code int) IError {
	e.code = code
	return e
}

func New(message string) (e IError) {
	e = &appError{code: http.StatusInternalServerError, Message: message}
	return
}
