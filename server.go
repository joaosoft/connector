package webserver

import (
	"fmt"
	"net"

	"time"

	"regexp"

	"strings"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type WebServer struct {
	config        *WebServerConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	routes        Routes
	middlewares   []MiddlewareFunc
	listener      net.Listener
	address       string
	errorhandler  ErrorHandler
}

func NewWebServer(options ...WebServerOption) (*WebServer, error) {
	pm := manager.NewManager(manager.WithRunInBackground(true), manager.WithLogLevel(logger.NoneLevel))
	log := logger.NewLogDefault("webserver", logger.WarnLevel)

	service := &WebServer{
		pm:          pm,
		logger:      log,
		routes:      make(Routes),
		middlewares: make([]MiddlewareFunc, 0),
		address:     ":80",
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		service.logger.Warn(err)
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.WebServer.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.WebServer
	if appConfig.WebServer.Address != "" {
		service.address = appConfig.WebServer.Address
	}

	service.Reconfigure(options...)

	service.AddRoute(MethodGet, "/favicon.ico", service.handlerFile)
	service.errorhandler = service.DefaultErrorHandler

	return service, nil
}

func (w *WebServer) AddMiddlewares(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *WebServer) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.routes[method] = append(w.routes[method], Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	})

	return nil
}

func (w *WebServer) AddRoutes(route ...Route) error {
	for _, r := range route {
		if err := w.AddRoute(r.Method, r.Path, r.Handler, r.Middlewares...); err != nil {
			return err
		}
	}
	return nil
}

func (w *WebServer) AddNamespace(path string, middlewares ...MiddlewareFunc) *Namespace {
	return &Namespace{
		Path:        path,
		Middlewares: middlewares,
		WebServer:   w,
	}
}

func (n *Namespace) AddRoutes(route ...Route) error {
	for _, r := range route {
		if err := n.WebServer.AddRoute(r.Method, fmt.Sprintf("%s%s", n.Path, r.Path), r.Handler, append(r.Middlewares, n.Middlewares...)...); err != nil {
			return err
		}
	}
	return nil
}

func (w *WebServer) SetErrorHandler(handler ErrorHandler) error {
	w.errorhandler = handler
	return nil
}

func (w *WebServer) Start() error {
	w.logger.Debug("executing Start")
	var err error

	w.listener, err = net.Listen("tcp", w.address)
	w.logger.Printf("http server started on %s", w.address)

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

func (w *WebServer) handleConnection(conn net.Conn) (err error) {
	var ctx *Context
	var handler HandlerFunc
	var length int
	startTime := time.Now()

	defer func() {
		conn.Close()
	}()

	// read response from connection
	request, err := NewRequest(conn)
	if err != nil {
		w.logger.Errorf("error getting request: [%s]", err)
		return err
	}

	// create response for request
	response := NewResponse(request)

	// create context with request and response
	ctx = NewContext(startTime, request, response)

	// middleware's of the server
	var route *Route
	route, err = w.GetRoute(request.Method, request.Url)
	if err != nil {
		w.logger.Errorf("error getting route: [%s]", err)
		goto on_error
	}

	// get url parameters
	if err = w.LoadUrlParms(request, route); err != nil {
		w.logger.Errorf("error loading url parameters: [%s]", err)
		goto on_error
	}

	// route handler
	handler = route.Handler

	// execute middlewares
	length = len(w.middlewares)
	for i, _ := range w.middlewares {
		if w.middlewares[length-1-i] != nil {
			handler = w.middlewares[length-1-i](handler)
		}
	}

	// middleware's of the specific route
	length = len(route.Middlewares)
	for i, _ := range route.Middlewares {
		if route.Middlewares[length-1-i] != nil {
			handler = route.Middlewares[length-1-i](handler)
		}
	}

	// run handlers with middleware's
	if err = handler(ctx); err != nil {
		w.logger.Errorf("error executing handler: [%s]", err)
		goto on_error
	}

on_error:
	if err != nil {
		w.errorhandler(ctx, err)
	}

	// write response to connection
	if err = ctx.Response.write(); err != nil {
		w.logger.Errorf("error writing response: [%s]", err)
	}

	w.logger.Printf("Address[ %s ] Method[ %s ] Url[ %s ] Protocol[ %s ] Elapsed[ %s ]", ctx.Request.IP, ctx.Request.Method, ctx.Request.Url, ctx.Request.Protocol, startTime)

	return nil
}

func (w *WebServer) Stop() error {
	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	return nil
}

func ConvertPathToRegex(path string) string {

	var re = regexp.MustCompile(`:[a-zA-Z0-9\-_]+[^/]`)
	regx := re.ReplaceAllString(path, `[a-zA-Z0-9-_]+[^/]`)

	return fmt.Sprintf("^%s", regx)
}

func (w *WebServer) GetRoute(method Method, url string) (*Route, error) {

	for _, route := range w.routes[method] {
		if regx, err := regexp.Compile(route.Regex); err != nil {
			return nil, err
		} else {
			if regx.MatchString(url) {
				return &route, nil
			}
		}
	}

	return nil, ErrorNotFound
}

func (w *WebServer) LoadUrlParms(request *Request, route *Route) error {

	routeUrl := strings.Split(route.Path, "/")
	url := strings.Split(request.Url, "/")

	for i, name := range routeUrl {
		if name != url[i] {
			request.UrlParms[name[1:]] = []string{url[i]}
		}
	}

	return nil
}