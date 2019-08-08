package connector

import (
	"errors"
)

var (
	ErrorNotFound = errors.New("route not found")
)

func (w *Server) DefaultErrorHandler(ctx *Context, err error) error {
	w.logger.Infof("handling error: %s", err)

	return ctx.Response.WithBody([]byte(err.Error()))
}
