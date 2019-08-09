package connector

import (
	"fmt"
)

type ClientManagerConfig struct {
	Services map[string]*ClientService `json:"services"`
	Log      Log                       `json:"log"`
}

type ClientService struct {
	Address string `json:"address"`
}

func NewClientManagerConfig() (*AppClientManagerConfig, error) {
	appConfig := &AppClientManagerConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
