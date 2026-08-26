package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseYAMLInvalidLineWrapsSentinel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("no_colon_here\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseYAML(path); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestLoadCheckedReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("no_colon_here\n"), 0644); err != nil {
		t.Fatal(err)
	}
	os.Setenv("CONTENT_MODERATION_CONFIG", path)
	defer os.Unsetenv("CONTENT_MODERATION_CONFIG")
	if _, err := LoadChecked(); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestParseYAMLScannerErrorWrapsSentinel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "long.yaml")
	content := "key:" + strings.Repeat("a", 70*1024) + "\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseYAML(path); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}
