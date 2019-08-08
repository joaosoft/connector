package connector

type AppConfig struct {
	Server ServerConfig `json:"server"`
	Client ClientConfig `json:"client"`
}
