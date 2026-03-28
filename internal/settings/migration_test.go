package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportTOML(t *testing.T) {
	dir := t.TempDir()
	tomlPath := filepath.Join(dir, "settings.toml")

	tomlContent := `[scooter]
brake-hibernation = true

[alarm]
enabled = true

[dashboard]
units = "metric"

[cellular]
apn = "internet"
`
	os.WriteFile(tomlPath, []byte(tomlContent), 0644)

	settings, err := ImportTOML(tomlPath)
	if err != nil {
		t.Fatalf("ImportTOML error: %v", err)
	}

	expected := map[string]string{
		"scooter.brake-hibernation": "true",
		"alarm.enabled":            "true",
		"dashboard.units":          "metric",
		"cellular.apn":             "internet",
	}

	for key, want := range expected {
		got, ok := settings[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if got != want {
			t.Errorf("settings[%q] = %q, want %q", key, got, want)
		}
	}
}

func TestImportTOML_FileNotFound(t *testing.T) {
	_, err := ImportTOML("/nonexistent/settings.toml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestImportTOML_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	tomlPath := filepath.Join(dir, "settings.toml")
	os.WriteFile(tomlPath, []byte(""), 0644)

	settings, err := ImportTOML(tomlPath)
	if err != nil {
		t.Fatalf("ImportTOML error: %v", err)
	}
	if len(settings) != 0 {
		t.Errorf("expected empty settings, got %d", len(settings))
	}
}
