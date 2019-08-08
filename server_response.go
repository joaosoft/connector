package connector

import (
	"io"
)

func (w *Server) NewResponse(request *Request) *Response {
	return &Response{
		Base:   request.Base,
		Writer: request.conn.(io.Writer),
	}
}
