package logger

import (
	"encoding/json"
	"io/ioutil"

	"go.uber.org/zap"
)

// New builds logger from json config
func New(path string) (*zap.SugaredLogger, error) {
	rawJSONConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := zap.Config{}
	if err := json.Unmarshal(rawJSONConfig, &config); err != nil {
		return nil, err
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
