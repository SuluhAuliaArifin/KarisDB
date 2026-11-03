package app

import (
	"runtime"

	"karisdb/internal/core"
)

const Version = "0.1.0"

type Status struct {
	Version string     `json:"version"`
	OS      string     `json:"os"`
	Arch    string     `json:"architecture"`
	Paths   core.Paths `json:"paths"`
}

var runtimeOS = runtime.GOOS
var runtimeArch = runtime.GOARCH
