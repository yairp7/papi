package config

type ServerConfig struct {
	Env             string `json:"env"`
	Version         string `json:"version"`
	Port            int64  `json:"port"`
	HeathCheckRoute string `json:"healthcheck_route"`
}
