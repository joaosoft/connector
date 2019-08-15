package connector

import (
	"io"
	"net"
)

func (c *Client) NewResponse(method string, address string, conn net.Conn) (*Response, error) {

	response := &Response{
		Base: Base{
			Client:  c,
			Method:  method,
			Address: address,
			Headers: make(Headers),
			conn:    conn,
		},
		Status: StatusOk,
		Reader: conn.(io.Reader),
	}

	return response, response.read()
}
