package config

type Config[T any] struct {
	Env             string `json:"env"`
	Version         string `json:"version"`
	Port            int64  `json:"port"`
	HeathCheckRoute string `json:"healthcheck_routh"`
	Data            T      `json:"data"`
}
