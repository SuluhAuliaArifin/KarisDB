package core

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestResolvePathsWindowsUsesLocalAppData(t *testing.T) {
	paths, err := resolvePaths("windows", func(string) string { return "" }, func() (string, error) {
		return `C:\Users\test`, nil
	}, `C:\Users\test\AppData\Local`)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(`C:\Users\test\AppData\Local`, "KarisDB"); paths.DataDir != want {
		t.Fatalf("DataDir = %q, want %q", paths.DataDir, want)
	}
	if want := filepath.Join(paths.DataDir, "data"); paths.DatabaseDir != want {
		t.Fatalf("DatabaseDir = %q, want %q", paths.DatabaseDir, want)
	}
}

func TestResolvePathsLinuxUsesXDGDataHome(t *testing.T) {
	paths, err := resolvePaths("linux", func(key string) string {
		if key == "XDG_DATA_HOME" {
			return "/tmp/user-data"
		}
		return ""
	}, func() (string, error) { return "/home/test", nil }, "")
	if err != nil {
		t.Fatal(err)
	}
	if want := "/tmp/user-data/karisdb"; paths.DataDir != want {
		t.Fatalf("DataDir = %q, want %q", paths.DataDir, want)
	}
}

func TestResolvePathsRequiresHomeWhenNeeded(t *testing.T) {
	_, err := resolvePaths("darwin", func(string) string { return "" }, func() (string, error) {
		return "", errors.New("home unavailable")
	}, "")
	if err == nil {
		t.Fatal("expected home directory error")
	}
}
