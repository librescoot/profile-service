package db

import (
	"testing"
)

func TestGlobalSettings_SetAndGet(t *testing.T) {
	store := openTestDB(t)

	err := store.SetGlobalSettings(map[string]string{
		"scooter.brake-hibernation": "true",
		"alarm.enabled":            "true",
	})
	if err != nil {
		t.Fatalf("SetGlobalSettings error: %v", err)
	}

	settings, err := store.GetGlobalSettings()
	if err != nil {
		t.Fatalf("GetGlobalSettings error: %v", err)
	}
	if settings["scooter.brake-hibernation"] != "true" {
		t.Errorf("brake-hibernation = %q, want %q", settings["scooter.brake-hibernation"], "true")
	}
	if settings["alarm.enabled"] != "true" {
		t.Errorf("alarm.enabled = %q, want %q", settings["alarm.enabled"], "true")
	}
}

func TestUserSettings_SetAndGet(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	err := store.SetUserSettings(p.ID, map[string]string{
		"dashboard.units": "metric",
	})
	if err != nil {
		t.Fatalf("SetUserSettings error: %v", err)
	}

	settings, err := store.GetUserSettings(p.ID)
	if err != nil {
		t.Fatalf("GetUserSettings error: %v", err)
	}
	if settings["dashboard.units"] != "metric" {
		t.Errorf("units = %q, want %q", settings["dashboard.units"], "metric")
	}
}

func TestUserSettings_DeletedWithProfile(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)
	store.SetUserSettings(p.ID, map[string]string{"dashboard.units": "metric"})

	store.DeleteProfile(p.ID)

	settings, err := store.GetUserSettings(p.ID)
	if err != nil {
		t.Fatalf("GetUserSettings error: %v", err)
	}
	if len(settings) != 0 {
		t.Errorf("expected empty settings after profile delete, got %d", len(settings))
	}
}

func TestSetGlobalSetting_Single(t *testing.T) {
	store := openTestDB(t)
	store.SetGlobalSettings(map[string]string{"alarm.enabled": "false"})

	err := store.SetGlobalSetting("alarm.enabled", "true")
	if err != nil {
		t.Fatalf("SetGlobalSetting error: %v", err)
	}

	settings, _ := store.GetGlobalSettings()
	if settings["alarm.enabled"] != "true" {
		t.Errorf("alarm.enabled = %q, want %q", settings["alarm.enabled"], "true")
	}
}

func TestSetUserSetting_Single(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	err := store.SetUserSetting(p.ID, "dashboard.units", "imperial")
	if err != nil {
		t.Fatalf("SetUserSetting error: %v", err)
	}

	settings, _ := store.GetUserSettings(p.ID)
	if settings["dashboard.units"] != "imperial" {
		t.Errorf("units = %q, want %q", settings["dashboard.units"], "imperial")
	}
}

func TestDeleteUserSetting(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)
	store.SetUserSetting(p.ID, "dashboard.units", "metric")

	err := store.DeleteUserSetting(p.ID, "dashboard.units")
	if err != nil {
		t.Fatalf("DeleteUserSetting error: %v", err)
	}

	settings, _ := store.GetUserSettings(p.ID)
	if _, exists := settings["dashboard.units"]; exists {
		t.Error("expected setting to be deleted")
	}
}
