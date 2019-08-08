package connector

import (
	"github.com/joaosoft/logger"
)

// ClientOption ...
type ClientOption func(client *Client)

// Reconfigure ...
func (c *Client) Reconfigure(options ...ClientOption) {
	for _, option := range options {
		option(c)
	}
}

// WithClientConfiguration ...
func WithClientConfiguration(config *ClientConfig) ClientOption {
	return func(client *Client) {
		client.config = config
	}
}

// WithClientLogger ...
func WithClientLogger(logger logger.ILogger) ClientOption {
	return func(client *Client) {
		client.logger = logger
		client.isLogExternal = true
	}
}

// WithClientLogLevel ...
func WithClientLogLevel(level logger.Level) ClientOption {
	return func(client *Client) {
		client.logger.SetLevel(level)
	}
}
