package connector

import (
	"fmt"
)

type ClientManagerConfig struct {
	Servers map[string]*ClientService `json:"Servers"`
	Log     Log                       `json:"log"`
}

type ClientService struct {
	Address string `json:"address"`
}

func NewClientManagerConfig() (*AppClientManagerConfig, error) {
	appConfig := &AppClientManagerConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
