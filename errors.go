package connector

import (
	"errors"
)

var (
	ErrorMethodNotFound        = errors.New("method not found")
	ErrorServiceNotFound       = errors.New("service not found")
	ErrorConfigurationNotFound = errors.New("configuration not found")
	ErrorServerDown            = errors.New("server is down")
)

func (w *Server) DefaultErrorHandler(ctx *Context, err error) error {
	w.logger.Infof("handling error: %s", err)

	ctx.Response.WithBody([]byte(err.Error())).WithStatus(StatusInternalError)

	return nil
}
