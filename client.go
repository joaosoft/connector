package connector

import (
	"fmt"
	"net"
	"time"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type Client struct {
	config        *ClientConfig
	isLogExternal bool
	logger        logger.ILogger
	dialer        net.Dialer
}

func NewClient(options ...ClientOption) (*Client, error) {
	config, err := NewClientConfig()

	service := &Client{
		logger: logger.NewLogDefault("client", logger.WarnLevel),
		config: &config.Client,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		level, _ := logger.ParseLevel(service.config.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	// create a new dialer to create connections
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	service.dialer = dialer

	service.Reconfigure(options...)

	return service, nil
}

func (c *Client) Invoke(method, address string, headers Headers, body ...[]byte) (*Response, error) {
	request := c.NewRequest(method, address, headers)

	if len(body) > 0 {
		request.WithBody(body[0])
	}

	return request.Send()
}


func (r *Request) Send() (*Response, error) {
	return r.Client.Send(r)
}

func (c *Client) Send(request *Request) (*Response, error) {
	startTime := time.Now()
	fmt.Println(color.WithColor("[IN] Method[%s] Address[%s] on Start[%s]", color.FormatBold, color.ForegroundBlue, color.BackgroundBlack, request.Method, request.Address, startTime))

	if c.logger.IsDebugEnabled() {
		if request.Body != nil {
			c.logger.Infof("[REQUEST BODY] [%s]", string(request.Body))
		}
	}

	c.logger.Debugf("executing method [%s] request to address [%s]", request.Method, request.Address)

	var conn net.Conn
	var err error

	conn, err = c.dialer.Dial("tcp", request.Address)

	if err != nil {
		return nil, err
	}

	body, err := request.build()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.Write(body)

	response, err := c.NewResponse(request.Method, request.Address, conn)

	if c.logger.IsDebugEnabled() {
		if response.Body != nil {
			c.logger.Infof("[RESPONSE BODY] [%s]", string(response.Body))
		}
	}

	fmt.Println(color.WithColor("[OUT] Method[%s] Address[%s] on Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, request.Method, request.Address, startTime, time.Since(startTime)))

	return response, err
}
