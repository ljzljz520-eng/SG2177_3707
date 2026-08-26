package config

import (
	"fmt"
	"os"
	"strings"
)

func (c Config) ValidateRuntime() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if strings.Contains(c.Listen, " ") {
		return fmt.Errorf("listen address contains spaces")
	}
	return nil
}

func (c Config) FileExists() bool {
	_, err := os.Stat(c.DBPath)
	return err == nil
}

func (c Config) Description() map[string]string {
	return map[string]string{"database": c.DBPath, "listen": c.Listen, "sort": c.SortMode, "backup": c.Backup, "http": fmt.Sprint(c.ServeHTTP)}
}

func (c Config) WithDatabase(path string) Config {
	c.DBPath = strings.TrimSpace(path)
	return c
}

func (c Config) WithSort(mode string) Config {
	c.SortMode = strings.TrimSpace(strings.ToLower(mode))
	return c
}

func (c Config) ShouldServe() bool {
	return c.ServeHTTP && c.Listen != ""
}
