package connector

import (
	"github.com/joaosoft/logger"
)

type ClientManager struct {
	config        *ClientManagerConfig
	isLogExternal bool
	logger        logger.ILogger
	client        *Client
}

func NewClientManager(options ...ClientManagerOption) (*ClientManager, error) {

	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	config, err := NewClientManagerConfig()

	service := &ClientManager{
		logger: logger.NewLogDefault("ClientManager", logger.LevelWarn),
		client: client,
		config: &config.ClientManager,
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

	service.Reconfigure(options...)

	if err = service.checkAlive(); err != nil {
		return nil, err
	}

	return service, nil
}

func (cm *ClientManager) checkAlive() error {
	for _, service := range cm.config.Servers {
		response, err := cm.client.NewRequest("alive", service.Address).Send()
		if err != nil {
			return err
		}

		if response.Status != StatusOk {
			return ErrorServerDown
		}
	}

	return nil
}

func (cm *ClientManager) Invoke(service, method string, headers Headers, body ...[]byte) (*Response, error) {
	var clientConf *ClientService
	var ok bool

	if clientConf, ok = cm.config.Servers[service]; !ok {
		return nil, ErrorConfigurationNotFound
	}

	return cm.client.Invoke(method, clientConf.Address, headers, body...)
}
