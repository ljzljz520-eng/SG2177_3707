package config

import "testing"

func TestConfigDefaults(t *testing.T) {
	config := Default()
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	if config.DBPath == "" || config.Listen == "" || config.SortMode != "number" {
		t.Fatalf("unexpected defaults: %#v", config)
	}
}

func TestConfigParse(t *testing.T) {
	config, err := Parse([]string{"-db", "tmp.db", "-sort", "date"})
	if err != nil || config.DBPath != "tmp.db" || config.SortMode != "date" {
		t.Fatalf("parse failed: %#v %v", config, err)
	}
}
