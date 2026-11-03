package app

import (
	"context"
	"encoding/json"
	"path/filepath"

	"karisdb/internal/core"
)

type Router struct {
	paths core.Paths
}

func NewRouter() (*Router, error) {
	paths, err := core.ResolvePaths()
	if err != nil {
		return nil, err
	}
	return &Router{paths: paths}, nil
}

func (router *Router) Execute(ctx context.Context, command string, arguments map[string]json.RawMessage) (any, *core.AppError) {
	if err := ctx.Err(); err != nil {
		return nil, core.NewError("REQUEST_CANCELED", "request was canceled")
	}
	if len(arguments) != 0 {
		return nil, core.NewError("INVALID_ARGUMENTS", "this command does not accept arguments")
	}

	switch command {
	case "help":
		return map[string]any{"commands": []string{"help", "version", "status"}}, nil
	case "version":
		return map[string]string{"version": Version}, nil
	case "status":
		return router.status()
	default:
		return nil, core.NewError("UNKNOWN_COMMAND", "unknown command")
	}
}

func (router *Router) status() (any, *core.AppError) {
	config, err := core.LoadConfig(filepath.Join(router.paths.ConfigDir, "config.json"))
	if err != nil {
		return nil, core.NewError("CONFIG_INVALID", "configuration could not be loaded")
	}
	paths := router.paths
	if config.DataDir != "" {
		paths = core.Paths{
			DataDir:     config.DataDir,
			RuntimeDir:  filepath.Join(config.DataDir, "runtimes"),
			DatabaseDir: filepath.Join(config.DataDir, "data"),
			ConfigDir:   filepath.Join(config.DataDir, "configs"),
			StateDir:    filepath.Join(config.DataDir, "state"),
			LogsDir:     filepath.Join(config.DataDir, "logs"),
		}
	}
	return Status{Version: Version, OS: runtimeOS, Arch: runtimeArch, Paths: paths}, nil
}
