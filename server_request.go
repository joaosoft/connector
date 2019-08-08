package connector

import (
	"io"
	"net"
)

func (w *Server) NewRequest(conn net.Conn, server *Server) (*Request, error) {

	request := &Request{
		Base: Base{
			Server:  server,
			IP:      conn.RemoteAddr().String(),
			Headers: make(Headers),
			conn:    conn,
		},
		Reader: conn.(io.Reader),
	}

	return request, request.read()
}
