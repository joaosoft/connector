package connector

import (
	"fmt"
)

type ServerManagerConfig struct {
	Log Log `json:"log"`
}

func NewServerManagerConfig() (*AppServerManagerConfig, error) {
	appConfig := &AppServerManagerConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
