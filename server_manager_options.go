package connector

import (
	"github.com/joaosoft/logger"
)

// ServerManagerOption ...
type ServerManagerOption func(ServerManager *ServerManager)

// Reconfigure ...
func (w *ServerManager) Reconfigure(options ...ServerManagerOption) {
	for _, option := range options {
		option(w)
	}
}

// WithServerManagerConfiguration ...
func WithServerManagerConfiguration(config *ServerManagerConfig) ServerManagerOption {
	return func(ServerManager *ServerManager) {
		ServerManager.config = config
	}
}

// WithServerManagerLogger ...
func WithServerManagerLogger(logger logger.ILogger) ServerManagerOption {
	return func(ServerManager *ServerManager) {
		ServerManager.logger = logger
		ServerManager.isLogExternal = true
	}
}

// WithServerManagerLogLevel ...
func WithServerManagerLogLevel(level logger.Level) ServerManagerOption {
	return func(ServerManager *ServerManager) {
		ServerManager.logger.SetLevel(level)
	}
}