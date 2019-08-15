package connector

import (
	"github.com/joaosoft/logger"
)

// ServerOption ...
type ServerOption func(server *Server)

// Reconfigure ...
func (w *Server) Reconfigure(options ...ServerOption) {
	for _, option := range options {
		option(w)
	}
}

// WithServerName ...
func WithServerName(name string) ServerOption {
	return func(server *Server) {
		server.name = name
	}
}

// WithServerConfiguration ...
func WithServerConfiguration(config *ServerConfig) ServerOption {
	return func(server *Server) {
		server.config = config
	}
}

// WithServerLogger ...
func WithServerLogger(logger logger.ILogger) ServerOption {
	return func(server *Server) {
		server.logger = logger
		server.isLogExternal = true
	}
}

// WithServerLogLevel ...
func WithServerLogLevel(level logger.Level) ServerOption {
	return func(server *Server) {
		server.logger.SetLevel(level)
	}
}

// WithServerAddress ...
func WithServerAddress(address string) ServerOption {
	return func(server *Server) {
		server.config.Address = address
	}
}
