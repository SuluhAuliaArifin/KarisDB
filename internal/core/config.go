package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const MaxConfigBytes = 1 << 20

type Config struct {
	DataDir string `json:"data_dir,omitempty"`
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	contents, err := io.ReadAll(io.LimitReader(file, MaxConfigBytes+1))
	if err != nil {
		return Config{}, err
	}
	if len(contents) > MaxConfigBytes {
		return Config{}, errors.New("configuration file exceeds size limit")
	}

	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return Config{}, errors.New("configuration must contain exactly one JSON object")
	}
	if config.DataDir != "" && !filepath.IsAbs(config.DataDir) {
		return Config{}, errors.New("data_dir must be an absolute path")
	}
	return config, nil
}
