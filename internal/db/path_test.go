package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDatabasePathNormalizesThingsDataDirectory(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "Things Database.thingsdatabase", thingsDatabaseFile)
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		t.Fatalf("create database bundle: %v", err)
	}
	if err := os.WriteFile(dbPath, []byte("sqlite"), 0o644); err != nil {
		t.Fatalf("write database file: %v", err)
	}

	got, err := ResolveDatabasePath(root)
	if err != nil {
		t.Fatalf("resolve database path: %v", err)
	}
	if got != dbPath {
		t.Fatalf("expected %q, got %q", dbPath, got)
	}
}

func TestResolveDatabasePathNormalizesThingsDatabaseDirectory(t *testing.T) {
	root := t.TempDir()
	bundle := filepath.Join(root, "Things Database.thingsdatabase")
	dbPath := filepath.Join(bundle, thingsDatabaseFile)
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatalf("create database bundle: %v", err)
	}
	if err := os.WriteFile(dbPath, []byte("sqlite"), 0o644); err != nil {
		t.Fatalf("write database file: %v", err)
	}

	got, err := ResolveDatabasePath(bundle)
	if err != nil {
		t.Fatalf("resolve database path: %v", err)
	}
	if got != dbPath {
		t.Fatalf("expected %q, got %q", dbPath, got)
	}
}

func TestResolveDatabasePathRejectsDirectoryWithoutDatabase(t *testing.T) {
	root := t.TempDir()

	got, err := ResolveDatabasePath(root)
	if err == nil {
		t.Fatalf("expected directory error, got path %q", got)
	}
	if !strings.Contains(err.Error(), "expected Things database file") {
		t.Fatalf("expected clear directory error, got %v", err)
	}
}
