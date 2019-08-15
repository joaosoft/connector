package connector

import (
	"io"
	"net"
	"time"
)

type Headers map[string][]string

type ErrorHandler func(ctx *Context, err error) error
type HandlerFunc func(ctx *Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Context struct {
	StartTime time.Time
	Request   *Request
	Response  *Response
}

type Base struct {
	IP      string
	Address string
	Method  string
	Headers Headers
	conn    net.Conn
	Server  *Server
	Client  *Client
}

type Request struct {
	Base
	Body   []byte
	Reader io.Reader
	Writer io.Writer
}

type Response struct {
	Base
	Body   []byte
	Status Status
	Reader io.Reader
	Writer io.Writer
}

type RequestHandler struct {
	Conn    net.Conn
	Handler HandlerFunc
}

type Servers map[string]*Server
type Clients map[string]*Client
