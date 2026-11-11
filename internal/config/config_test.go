package config

import "testing"

func TestLoadDefaultsAndRequiredDatabase(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/tasks")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("API_KEY", "secret")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.HTTPAddr != ":8080" || c.APIKey != "secret" {
		t.Fatalf("unexpected config: %+v", c)
	}
}

func TestLoadRequiresDatabase(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("expected missing DATABASE_URL error")
	}
}
