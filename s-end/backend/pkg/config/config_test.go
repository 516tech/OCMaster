package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := Load()
	if cfg.ServerPort == "" {
		t.Error("expected default ServerPort")
	}
	if cfg.DBType != "sqlite" {
		t.Errorf("expected sqlite, got %s", cfg.DBType)
	}
	if cfg.DBPath == "" {
		t.Error("expected default DBPath")
	}
	if cfg.JWTSecret == "" {
		t.Error("expected default JWTSecret")
	}
}

func TestLoad_DefaultsAreSane(t *testing.T) {
	cfg := Load()
	if cfg.ServerPort != "8080" {
		t.Errorf("expected port 8080, got %s", cfg.ServerPort)
	}
}

func TestOpenDB_SQLite_Memory(t *testing.T) {
	cfg := &Config{DBType: "sqlite", DBPath: ":memory:"}
	db, err := cfg.OpenDB()
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func TestOpenDB_SQLite_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	cfg := &Config{DBType: "sqlite", DBPath: path}
	db, err := cfg.OpenDB()
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected db file to exist")
	}
}
