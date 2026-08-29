package simpleconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		DSN string
	}
}

func TestLoad_DefaultPaths(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	os.Mkdir("config", 0755)
	os.WriteFile("config/config.toml", []byte(`
[server]
host = "localhost"
port = 8080
`), 0644)

	res, err := Load[testConfig]()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if res.Config.Server.Host != "localhost" {
		t.Errorf("expected host=localhost, got %q", res.Config.Server.Host)
	}
	if res.Config.Server.Port != 8080 {
		t.Errorf("expected port=8080, got %d", res.Config.Server.Port)
	}
	if res.LoadedPath != filepath.Join(".", "config", "config.toml") {
		t.Errorf("expected loaded path %q, got %q", filepath.Join(".", "config", "config.toml"), res.LoadedPath)
	}
}

func TestLoad_CustomFilename(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	os.WriteFile("myconfig.toml", []byte(`
[server]
host = "custom"
port = 9090
`), 0644)

	res, err := Load[testConfig](WithFilename("myconfig.toml"))
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if res.Config.Server.Host != "custom" {
		t.Errorf("expected host=custom, got %q", res.Config.Server.Host)
	}
}

func TestLoad_CustomPaths(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	customDir := filepath.Join(tmpDir, "etc", "myapp")
	os.MkdirAll(customDir, 0755)
	os.WriteFile(filepath.Join(customDir, "config.toml"), []byte(`
[server]
host = "frometc"
port = 7070
`), 0644)

	res, err := Load[testConfig](WithCustomPaths(customDir))
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if res.Config.Server.Host != "frometc" {
		t.Errorf("expected host=frometc, got %q", res.Config.Server.Host)
	}
}

func TestLoad_CustomPathsPriority(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	os.Mkdir("config", 0755)
	os.WriteFile("config/config.toml", []byte(`
[server]
host = "local"
port = 8080
`), 0644)

	customDir := filepath.Join(tmpDir, "etc", "myapp")
	os.MkdirAll(customDir, 0755)
	os.WriteFile(filepath.Join(customDir, "config.toml"), []byte(`
[server]
host = "custom"
port = 9090
`), 0644)

	res, err := Load[testConfig](WithCustomPaths(customDir))
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if res.Config.Server.Host != "custom" {
		t.Errorf("custom path should have priority, got %q", res.Config.Server.Host)
	}
}

func TestLoad_Validator(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	os.Mkdir("config", 0755)
	os.WriteFile("config/config.toml", []byte(`
[server]
host = ""
port = 8080
`), 0644)

	validator := func(cfg any) error {
		c := cfg.(*testConfig)
		if c.Server.Host == "" {
			return fmt.Errorf("server.host is required")
		}
		return nil
	}

	_, err := Load[testConfig](WithValidator(validator))
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoad_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	_, err := Load[testConfig]()
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestMustLoad(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	os.Mkdir("config", 0755)
	os.WriteFile("config/config.toml", []byte(`
[server]
host = "mustload"
port = 8080
`), 0644)

	res := MustLoad[testConfig]()
	if res.Config.Server.Host != "mustload" {
		t.Errorf("expected host=mustload, got %q", res.Config.Server.Host)
	}
}
