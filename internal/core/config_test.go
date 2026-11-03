package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAcceptsAbsoluteDataDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"data_dir":"C:\\KarisDBData"}`), 0600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.DataDir != `C:\KarisDBData` {
		t.Fatalf("DataDir = %q", config.DataDir)
	}
}

func TestLoadConfigRejectsRelativeDataDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"data_dir":"relative"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected relative data_dir to be rejected")
	}
}

func TestLoadConfigRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"unexpected":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected unknown config field to be rejected")
	}
}
