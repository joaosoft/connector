package connector

type AppServerConfig struct {
	Server ServerConfig `json:"server"`
}

type AppClientConfig struct {
	Client ClientConfig `json:"client"`
}

type AppServerManagerConfig struct {
	ServerManager ServerManagerConfig `json:"server_manager"`
}

type AppClientManagerConfig struct {
	ClientManager ClientManagerConfig `json:"client_manager"`
}

type Log struct {
	Level string `json:"level"`
}
