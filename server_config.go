package connector

import (
	"fmt"
)

type ServerConfig struct {
	Address string `json:"address"`
	Log     Log    `json:"log"`
}

func NewServerConfig() (*AppServerConfig, error) {
	appConfig := &AppServerConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
