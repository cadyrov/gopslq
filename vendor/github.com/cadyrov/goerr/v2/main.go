package goerr

import (
	"net/http"
)

type IError interface {
	Error() string
	Code() int
	Details() []IError
	PushDetail(iError IError)
	Tag(string)
	GetTag() string
	GetError() error
}

type baseError struct {
	err     error
	code    int
	tag     string
	details []IError
}

func New(code int, err error) IError {
	return &baseError{
		code: code,
		err:  err,
	}
}

func Unauthorized(err error) IError {
	return New(http.StatusUnauthorized, err)
}

func Forbidden(err error) IError {
	return New(http.StatusForbidden, err)
}

func BadRequest(err error) IError {
	return New(http.StatusBadRequest, err)
}

func NotFound(err error) IError {
	return New(http.StatusNotFound, err)
}

func NotAllowed(err error) IError {
	return New(http.StatusMethodNotAllowed, err)
}

func NotAcceptable(err error) IError {
	return New(http.StatusNotAcceptable, err)
}

func Conflict(err error) IError {
	return New(http.StatusConflict, err)
}

func Unprocessable(err error) IError {
	return New(http.StatusUnprocessableEntity, err)
}

func Internal(err error) IError {
	return New(http.StatusInternalServerError, err)
}

func (c *baseError) Tag(tag string) {
	c.tag = tag
}

func (c *baseError) GetTag() string {
	return c.tag
}

func (c baseError) GetError() error {
	return c.err
}

func (c *baseError) PushDetail(iError IError) {
	c.details = append(c.details, &baseError{
		code:    iError.Code(),
		err:     iError.GetError(),
		details: iError.Details(),
		tag:     iError.GetTag(),
	})
}

func (c baseError) Error() string {
	return c.err.Error()
}

func (c baseError) Code() int {
	return c.code
}

func (c baseError) Details() []IError {
	if len(c.details) == 0 {
		return nil
	}

	res := make([]IError, 0, len(c.details))

	for i := range c.details {
		res = append(res, &baseError{
			code:    c.details[i].Code(),
			err:     c.GetError(),
			details: c.details[i].Details(),
		})
	}

	return res
}