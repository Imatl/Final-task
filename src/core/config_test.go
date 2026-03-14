package core

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.properties")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadConfig(t *testing.T) {
	path := writeTestConfig(t, `
service.name=testapp
service.port=9090
# comment line
empty.key=

db.host=localhost
`)
	if err := LoadConfig(path); err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if got := GetString("service.name", ""); got != "testapp" {
		t.Errorf("service.name = %q, want %q", got, "testapp")
	}
	if got := GetInt("service.port", 0); got != 9090 {
		t.Errorf("service.port = %d, want %d", got, 9090)
	}
	if got := GetString("db.host", ""); got != "localhost" {
		t.Errorf("db.host = %q, want %q", got, "localhost")
	}
}

func TestGetStringDefault(t *testing.T) {
	path := writeTestConfig(t, "key=value")
	LoadConfig(path)

	if got := GetString("missing", "fallback"); got != "fallback" {
		t.Errorf("missing key = %q, want %q", got, "fallback")
	}
}

func TestGetStringFromVault(t *testing.T) {
	path := writeTestConfig(t, "secret=FROM_VAULT")
	LoadConfig(path)

	if got := GetString("secret", "default"); got != "default" {
		t.Errorf("FROM_VAULT should return default, got %q", got)
	}
}

func TestGetIntDefault(t *testing.T) {
	path := writeTestConfig(t, "key=notanumber")
	LoadConfig(path)

	if got := GetInt("key", 42); got != 42 {
		t.Errorf("invalid int = %d, want %d", got, 42)
	}
	if got := GetInt("missing", 99); got != 99 {
		t.Errorf("missing int = %d, want %d", got, 99)
	}
}

func TestGetBool(t *testing.T) {
	path := writeTestConfig(t, "flag=true\nflag2=false\nflag3=invalid")
	LoadConfig(path)

	if got := GetBool("flag", false); got != true {
		t.Errorf("flag = %v, want true", got)
	}
	if got := GetBool("flag2", true); got != false {
		t.Errorf("flag2 = %v, want false", got)
	}
	if got := GetBool("flag3", true); got != true {
		t.Errorf("flag3 invalid should return default true, got %v", got)
	}
}

func TestEnvOverride(t *testing.T) {
	os.Setenv("DB_HOST", "envhost")
	defer os.Unsetenv("DB_HOST")

	path := writeTestConfig(t, "db.host=filehost")
	LoadConfig(path)

	if got := GetString("db.host", ""); got != "envhost" {
		t.Errorf("env override = %q, want %q", got, "envhost")
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	if err := LoadConfig("/nonexistent/path.properties"); err == nil {
		t.Error("expected error for missing config file")
	}
}

func TestCommentAndEmptyLines(t *testing.T) {
	path := writeTestConfig(t, `
# full comment
   # indented comment

key1=val1
  key2 = val2
`)
	LoadConfig(path)

	if got := GetString("key1", ""); got != "val1" {
		t.Errorf("key1 = %q, want %q", got, "val1")
	}
	if got := GetString("key2", ""); got != "val2" {
		t.Errorf("key2 = %q, want %q", got, "val2")
	}
}
