package connector

import (
	"github.com/joaosoft/logger"
)

// ClientManagerOption ...
type ClientManagerOption func(ClientManager *ClientManager)

// Reconfigure ...
func (c *ClientManager) Reconfigure(options ...ClientManagerOption) {
	for _, option := range options {
		option(c)
	}
}

// WithClientManagerConfiguration ...
func WithClientManagerConfiguration(config *ClientManagerConfig) ClientManagerOption {
	return func(ClientManager *ClientManager) {
		ClientManager.config = config
	}
}

// WithClientManagerLogger ...
func WithClientManagerLogger(logger logger.ILogger) ClientManagerOption {
	return func(ClientManager *ClientManager) {
		ClientManager.logger = logger
		ClientManager.isLogExternal = true
	}
}

// WithClientManagerLogLevel ...
func WithClientManagerLogLevel(level logger.Level) ClientManagerOption {
	return func(ClientManager *ClientManager) {
		ClientManager.logger.SetLevel(level)
	}
}
