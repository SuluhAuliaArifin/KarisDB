package core

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type Paths struct {
	DataDir     string `json:"data_dir"`
	RuntimeDir  string `json:"runtime_dir"`
	DatabaseDir string `json:"database_dir"`
	ConfigDir   string `json:"config_dir"`
	StateDir    string `json:"state_dir"`
	LogsDir     string `json:"logs_dir"`
}

func ResolvePaths() (Paths, error) {
	return resolvePaths(runtime.GOOS, os.Getenv, os.UserHomeDir, os.Getenv("LOCALAPPDATA"))
}

func resolvePaths(goos string, getenv func(string) string, homeDir func() (string, error), localAppData string) (Paths, error) {
	var dataDir string
	joinPath := path.Join
	isAbsPath := path.IsAbs

	switch goos {
	case "windows":
		joinPath = filepath.Join
		isAbsPath = filepath.IsAbs
		root := localAppData
		if root == "" {
			home, err := homeDir()
			if err != nil {
				return Paths{}, err
			}
			root = joinPath(home, "AppData", "Local")
		}
		dataDir = joinPath(root, "KarisDB")
	case "darwin":
		home, err := homeDir()
		if err != nil {
			return Paths{}, err
		}
		dataDir = joinPath(home, "Library", "Application Support", "KarisDB")
	default:
		root := getenv("XDG_DATA_HOME")
		if root == "" || !isAbsPath(root) {
			home, err := homeDir()
			if err != nil {
				return Paths{}, err
			}
			root = joinPath(home, ".local", "share")
		}
		dataDir = joinPath(root, "karisdb")
	}

	if dataDir == "" {
		return Paths{}, errors.New("could not resolve KarisDB data directory")
	}

	return Paths{
		DataDir:     dataDir,
		RuntimeDir:  joinPath(dataDir, "runtimes"),
		DatabaseDir: joinPath(dataDir, "data"),
		ConfigDir:   joinPath(dataDir, "configs"),
		StateDir:    joinPath(dataDir, "state"),
		LogsDir:     joinPath(dataDir, "logs"),
	}, nil
}
