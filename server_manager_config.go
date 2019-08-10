package connector

import (
	"fmt"
)

type ServerManagerConfig struct {
	Services map[string]*ServerService `json:"services"`
	Log Log `json:"log"`
}

type ServerService struct {
	Address string `json:"address"`
}

func NewServerManagerConfig() (*AppServerManagerConfig, error) {
	appConfig := &AppServerManagerConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
