package server

import (
	"fmt"
	"web/common"
)

var (
	ErrorNotFound = NewError(common.StatusNotFound, "route not found")
)

type Error struct {
	Status   common.Status
	Messages interface{}
}

func NewError(status common.Status, errors ...interface{}) *Error {
	err := &Error{
		Status: status,
	}

	if len(errors) > 0 {
		err.Messages = errors
	} else {
		err.Messages = []string{common.StatusText(status)}
	}

	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf("status=%d, messages=%v", e.Status, e.Messages)
}

func (w *Server) DefaultErrorHandler(ctx *Context, err error) error {
	w.logger.Infof("handling error: %s", err)

	if e, ok := err.(*Error); ok {
		return ctx.Response.JSON(e.Status, e)
	}

	return ctx.Response.JSON(common.StatusInternalServerError, NewError(common.StatusInternalServerError, err))
}
