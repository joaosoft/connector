package connector

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type ServerManager struct {
	config        *ServerManagerConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	servers       Servers
	started       bool
}

func NewServerManager(options ...ServerManagerOption) (*ServerManager, error) {
	config, err := NewServerManagerConfig()
	pm := manager.NewManager()

	service := &ServerManager{
		logger:  logger.NewLogDefault("ServerManager", logger.WarnLevel),
		servers: make(Servers),
		pm:      pm,
		config:  &config.ServerManager,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		level, _ := logger.ParseLevel(service.config.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	return service, nil
}

func (sm *ServerManager) Register(service string, server *Server) *ServerManager {
	server.AddMethod("alive", handleAlive)

	if config, ok := sm.config.Services[service]; ok {
		server.address = config.Address
	}

	sm.servers[service] = server
	return sm
}

func (sm *ServerManager) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	for service, server := range sm.servers {
		sm.logger.Infof("adding server [service: %s]", service)
		sm.pm.AddProcess(service, server)
	}

	sm.started = true
	wg.Done()

	return sm.pm.Start()
}

func (sm *ServerManager) Started() bool {
	return sm.started
}

func (sm *ServerManager) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}
	defer wg.Done()

	sm.started = false

	return sm.pm.Stop()
}

func handleAlive(ctx *Context) error {
	data := struct {
		Message string    `json:"message"`
		Time    time.Time `json:"time"`
	}{
		Message: "I'm alive!",
		Time:    time.Now().UTC(),
	}

	byts, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx.Response.WithBody(byts).withStatus(StatusOk)

	return nil
}
