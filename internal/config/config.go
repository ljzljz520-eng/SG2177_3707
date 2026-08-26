package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DBPath    string
	Listen    string
	SortMode  string
	Backup    string
	ServeHTTP bool
}

func Default() Config {
	return Config{DBPath: "equipment.db", Listen: "127.0.0.1:8080", SortMode: "number", Backup: "equipment-backup.json"}
}

func FromEnv(base Config) Config {
	if value := strings.TrimSpace(os.Getenv("EQUIPMENT_DB")); value != "" {
		base.DBPath = value
	}
	if value := strings.TrimSpace(os.Getenv("EQUIPMENT_LISTEN")); value != "" {
		base.Listen = value
	}
	if value := strings.TrimSpace(os.Getenv("EQUIPMENT_SORT")); value != "" {
		base.SortMode = value
	}
	if value := strings.TrimSpace(os.Getenv("EQUIPMENT_BACKUP")); value != "" {
		base.Backup = value
	}
	if strings.EqualFold(os.Getenv("EQUIPMENT_HTTP"), "true") {
		base.ServeHTTP = true
	}
	return base
}

func Parse(args []string) (Config, error) {
	base := FromEnv(Default())
	flags := flag.NewFlagSet("equipment-service", flag.ContinueOnError)
	dbPath := flags.String("db", base.DBPath, "SQLite database path")
	listen := flags.String("listen", base.Listen, "HTTP listen address")
	sortMode := flags.String("sort", base.SortMode, "default sort mode")
	backup := flags.String("backup", base.Backup, "backup file path")
	httpMode := flags.Bool("http", base.ServeHTTP, "serve HTTP")
	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	config := Config{DBPath: strings.TrimSpace(*dbPath), Listen: strings.TrimSpace(*listen), SortMode: strings.TrimSpace(*sortMode), Backup: strings.TrimSpace(*backup), ServeHTTP: *httpMode}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c Config) Validate() error {
	if c.DBPath == "" {
		return errors.New("database path is required")
	}
	if c.Listen == "" {
		return errors.New("listen address is required")
	}
	if c.SortMode != "number" && c.SortMode != "date" {
		return fmt.Errorf("unsupported sort mode %q", c.SortMode)
	}
	return nil
}

func (c Config) EnsureParent() error {
	parent := filepath.Dir(c.DBPath)
	if parent == "." || parent == "" {
		return nil
	}
	return os.MkdirAll(parent, 0750)
}

func (c Config) Summary() string {
	return fmt.Sprintf("db=%s listen=%s sort=%s http=%t", c.DBPath, c.Listen, c.SortMode, c.ServeHTTP)
}
