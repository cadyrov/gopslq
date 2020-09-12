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
	HTTP(code int) IError
}

type AppError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Detail  []IError `json:"detail,omitempty"`
}

func (e *AppError) PushDetail(ae IError) {
	e.Detail = append(e.Detail, ae)
}

func (e *AppError) Error() (er string) {
	er += fmt.Sprintf("Code: %v; ", e.Code)

	er += "Msg: " + e.Message + ";  "

	if len(e.GetDetails()) == 0 {
		return
	}

	er += " Details: {"

	for idx := range e.GetDetails() {
		er += e.GetDetails()[idx].Error()
	}

	er += "}"

	return er
}

func (e *AppError) GetCode() int {
	return e.Code
}

func (e *AppError) GetMessage() string {
	return e.Message
}

func (e *AppError) GetDetails() []IError {
	return e.Detail
}

func (e *AppError) HTTP(code int) IError {
	e.Code = code

	return e
}

func New(message string) IError {
	e := &AppError{Code: http.StatusInternalServerError, Message: message}

	return e
}
