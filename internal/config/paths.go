package config

import (
	"path/filepath"
	"strings"
)

func ResolvePath(base, child string) string {
	if strings.TrimSpace(child) == "" {
		return base
	}
	if filepath.IsAbs(child) {
		return filepath.Clean(child)
	}
	return filepath.Join(base, child)
}

func IsMemoryDatabase(path string) bool {
	clean := strings.ToLower(strings.TrimSpace(path))
	return clean == ":memory:" || clean == "file::memory:?cache=shared"
}

func BackupPath(dbPath string) string {
	if IsMemoryDatabase(dbPath) {
		return "equipment-backup.json"
	}
	return dbPath + ".json"
}

func WithDefaults(config Config) Config {
	defaults := Default()
	if config.DBPath == "" {
		config.DBPath = defaults.DBPath
	}
	if config.Listen == "" {
		config.Listen = defaults.Listen
	}
	if config.SortMode == "" {
		config.SortMode = defaults.SortMode
	}
	if config.Backup == "" {
		config.Backup = defaults.Backup
	}
	return config
}
