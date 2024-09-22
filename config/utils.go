package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

func ReadBase64Config[T any](b64 string) (Config[T], error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Config[T]{}, err
	}

	return readConfig[T](data)
}

func ReadConfig[T any](path string) (Config[T], error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config[T]{}, err
	}

	return readConfig[T](data)
}

func readConfig[T any](configBytes []byte) (Config[T], error) {
	var config Config[T]
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		return Config[T]{}, err
	}

	return config, nil
}
