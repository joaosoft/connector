package connector

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type Server struct {
	config        *ServerConfig
	isLogExternal bool
	logger        logger.ILogger
	methods       Methods
	middlewares   []MiddlewareFunc
	listener      net.Listener
	address       string
	errorhandler  ErrorHandler
}

func NewServer(options ...ServerOption) (*Server, error) {
	config, err := NewServerConfig()

	service := &Server{
		logger:      logger.NewLogDefault("server", logger.WarnLevel),
		methods:     make(Methods),
		middlewares: make([]MiddlewareFunc, 0),
		address:     ":5000",
		config:      &ServerConfig{},
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		level, _ := logger.ParseLevel(config.Client.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	if config.Server.Address != "" {
		service.address = config.Server.Address
	}

	service.errorhandler = service.DefaultErrorHandler

	return service, nil
}

func (w *Server) AddMiddlewares(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *Server) AddMethod(method string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.methods[method] = NewMethod(method, handler, middleware...)

	return nil
}

func (w *Server) AddMethods(methods ...*Method) error {
	for _, r := range methods {
		if err := w.AddMethod(r.Method, r.Handler, r.Middlewares...); err != nil {
			return err
		}
	}
	return nil
}

func (w *Server) SetErrorHandler(handler ErrorHandler) error {
	w.errorhandler = handler
	return nil
}

func (w *Server) Start() error {
	w.logger.Debug("executing Start")
	var err error

	w.listener, err = net.Listen("tcp", w.address)
	if err != nil {
		w.logger.Errorf("error connecting to %s: %s", w.address, err)
		return err
	}

	if w.address == ":0" {
		split := strings.Split(w.listener.Addr().String(), ":")
		w.address = fmt.Sprintf(":%s", split[len(split)-1])
	}

	fmt.Println(color.WithColor("Connector server started on [%s]", color.FormatBold, color.ForegroundRed, color.BackgroundBlack, w.address))

	for {
		conn, err := w.listener.Accept()
		w.logger.Info("accepted connection")
		if err != nil {
			w.logger.Errorf("error accepting connection: %s", err)
			continue
		}

		if conn == nil {
			w.logger.Error("the connection isn't initialized")
			continue
		}

		go w.handleConnection(conn)
	}

	return err
}

func (w *Server) handleConnection(conn net.Conn) (err error) {
	var ctx *Context
	var length int
	var handlerRoute HandlerFunc
	startTime := time.Now()

	defer func() {
		conn.Close()
	}()

	// read response from connection
	request, err := w.NewRequest(conn, w)
	if err != nil {
		w.logger.Errorf("error getting request: [%s]", err)
		return err
	}

	fmt.Println(color.WithColor("[IN] IP[%s] Method[%s] Start[%s]", color.FormatBold, color.ForegroundBlue, color.BackgroundBlack, request.IP, request.Method, startTime))

	if w.logger.IsDebugEnabled() {
		if request.Body != nil {
			w.logger.Infof("[REQUEST BODY] [%s]", string(request.Body))
		}
	}

	// create response for request
	response := w.NewResponse(request)

	// create context with request and response
	ctx = NewContext(startTime, request, response)

	// route of the Server
	method, err := w.GetMethod(request.Method)
	if err != nil {
		w.logger.Errorf("error getting route: [%s]", err)
		goto done
	}

	// route handler
	handlerRoute = emptyHandler
	handlerRoute = func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {
			if err := method.Handler(ctx); err != nil {
				return err
			}

			return next(ctx)
		}

	}(handlerRoute)

	// execute middlewares
	length = len(w.middlewares)
	for i, _ := range w.middlewares {
		if w.middlewares[length-1-i] != nil {
			handlerRoute = w.middlewares[length-1-i](handlerRoute)
		}
	}

	// middleware's of the specific method
	length = len(method.Middlewares)
	for i, _ := range method.Middlewares {
		if method.Middlewares[length-1-i] != nil {
			handlerRoute = method.Middlewares[length-1-i](handlerRoute)
		}
	}

	// run handlers with middleware's
	if err = handlerRoute(ctx); err != nil {
		w.logger.Errorf("error executing handler: [%s]", err)
		goto done
	}

done:
	if err != nil {
		if er := w.errorhandler(ctx, err); er != nil {
			w.logger.Errorf("error writing error: [error: %s] %s", err, er)
		}
	}

	// write response to connection
	if err = ctx.Response.write(); err != nil {
		w.logger.Errorf("error writing response: [%s]", err)
	}

	fmt.Println(color.WithColor("[OUT] IP[%s] Method[%s] Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, ctx.Request.IP, ctx.Request.Method, startTime, time.Since(startTime)))

	return nil
}

func (w *Server) Stop() error {
	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	return nil
}

func (w *Server) GetMethod(method string) (*Method, error) {
	if m, ok := w.methods[method]; ok {
		return m, nil
	}

	return nil, ErrorNotFound
}

func emptyHandler(ctx *Context) error {
	return nil
}
