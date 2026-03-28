package settings

import (
	"testing"
)

func TestMerge_GlobalOnly(t *testing.T) {
	global := map[string]string{
		"scooter.brake-hibernation": "true",
		"alarm.enabled":            "true",
	}

	result := Merge(global, nil)
	if result["scooter.brake-hibernation"] != "true" {
		t.Errorf("expected global value preserved")
	}
}

func TestMerge_UserOverridesAllowedCategory(t *testing.T) {
	global := map[string]string{
		"dashboard.units":           "metric",
		"scooter.brake-hibernation": "true",
	}
	user := map[string]string{
		"dashboard.units": "imperial",
	}

	result := Merge(global, user)
	if result["dashboard.units"] != "imperial" {
		t.Errorf("dashboard.units = %q, want %q", result["dashboard.units"], "imperial")
	}
	if result["scooter.brake-hibernation"] != "true" {
		t.Errorf("global setting should be preserved")
	}
}

func TestMerge_UserCannotOverrideGlobalCategory(t *testing.T) {
	global := map[string]string{
		"alarm.enabled": "true",
	}
	user := map[string]string{
		"alarm.enabled": "false",
	}

	result := Merge(global, user)
	if result["alarm.enabled"] != "true" {
		t.Errorf("alarm.enabled = %q, want %q (global should win)", result["alarm.enabled"], "true")
	}
}

func TestIsUserOverridable(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"dashboard.units", true},
		{"dashboard.brightness", true},
		{"navigation.routing-mode", true},
		{"scooter.brake-hibernation", false},
		{"alarm.enabled", false},
		{"cellular.apn", false},
		{"engine-ecu.kers", false},
		{"updates.auto-install", false},
		{"keycard.master-uid", false},
		{"unknown.setting", false},
	}

	for _, tt := range tests {
		got := IsUserOverridable(tt.key)
		if got != tt.want {
			t.Errorf("IsUserOverridable(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}
