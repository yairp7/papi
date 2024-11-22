package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

func ReadBase64Config[T any](b64 string) (*T, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	return readConfig[T](data)
}

func ReadConfig[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return readConfig[T](data)
}

func readConfig[T any](configBytes []byte) (*T, error) {
	var config T
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
