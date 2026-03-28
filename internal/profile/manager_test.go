package profile

import (
	"testing"
)

type fakePublisher struct {
	profileFields  map[string]string
	settingsFields map[string]string
}

func (f *fakePublisher) PublishProfile(fields map[string]string) {
	f.profileFields = fields
}

func (f *fakePublisher) PublishSettings(fields map[string]string) {
	f.settingsFields = fields
}

func (f *fakePublisher) ClearProfile() {
	f.profileFields = map[string]string{"active.id": ""}
}

func TestManager_Activate(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	store.SetGlobalSettings(map[string]string{
		"alarm.enabled":   "true",
		"dashboard.units": "metric",
	})

	p, _ := store.CreateProfile("Test User", "user", 42)
	store.SetUserSettings(p.ID, map[string]string{
		"dashboard.units": "imperial",
	})

	err := mgr.Activate(p)
	if err != nil {
		t.Fatalf("Activate error: %v", err)
	}

	if pub.profileFields["active.id"] != p.ID {
		t.Errorf("active.id = %q, want %q", pub.profileFields["active.id"], p.ID)
	}
	if pub.settingsFields["dashboard.units"] != "imperial" {
		t.Errorf("dashboard.units = %q, want %q", pub.settingsFields["dashboard.units"], "imperial")
	}
	if pub.settingsFields["alarm.enabled"] != "true" {
		t.Errorf("alarm.enabled = %q, want %q", pub.settingsFields["alarm.enabled"], "true")
	}
}

func TestManager_Deactivate(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	store.SetGlobalSettings(map[string]string{
		"alarm.enabled":   "true",
		"dashboard.units": "metric",
	})

	p, _ := store.CreateProfile("Test", "user", 1)
	mgr.Activate(p)

	err := mgr.Deactivate()
	if err != nil {
		t.Fatalf("Deactivate error: %v", err)
	}

	if pub.profileFields["active.id"] != "" {
		t.Errorf("active.id = %q, want empty", pub.profileFields["active.id"])
	}
	if pub.settingsFields["dashboard.units"] != "metric" {
		t.Errorf("dashboard.units = %q, want %q (reverted to global)", pub.settingsFields["dashboard.units"], "metric")
	}
}

func TestManager_ActiveProfile(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	if mgr.ActiveProfile() != nil {
		t.Error("expected nil active profile initially")
	}

	p, _ := store.CreateProfile("Test", "user", 1)
	mgr.Activate(p)

	if mgr.ActiveProfile() == nil {
		t.Error("expected non-nil active profile after activation")
	}
	if mgr.ActiveProfile().ID != p.ID {
		t.Errorf("active ID = %q, want %q", mgr.ActiveProfile().ID, p.ID)
	}
}

func TestManager_HandleSettingChange_OwnerPersistsGlobal(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	store.SetGlobalSettings(map[string]string{"alarm.enabled": "true"})
	owner, _ := store.CreateProfile("Owner", "owner", 1)
	mgr.Activate(owner)

	mgr.HandleSettingChange("alarm.enabled", "false")

	global, _ := store.GetGlobalSettings()
	if global["alarm.enabled"] != "false" {
		t.Errorf("alarm.enabled = %q, want %q (owner should persist global)", global["alarm.enabled"], "false")
	}
}

func TestManager_HandleSettingChange_UserPersistsOverridable(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	store.SetGlobalSettings(map[string]string{"dashboard.units": "metric"})
	user, _ := store.CreateProfile("User", "user", 1)
	mgr.Activate(user)

	mgr.HandleSettingChange("dashboard.units", "imperial")

	userSettings, _ := store.GetUserSettings(user.ID)
	if userSettings["dashboard.units"] != "imperial" {
		t.Errorf("user setting = %q, want %q", userSettings["dashboard.units"], "imperial")
	}
}

func TestManager_HandleSettingChange_GuestDoesNotPersist(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	store.SetGlobalSettings(map[string]string{"dashboard.units": "metric"})
	guest, _ := store.CreateProfile("Guest", "guest", 1)
	mgr.Activate(guest)

	mgr.HandleSettingChange("dashboard.units", "imperial")

	userSettings, _ := store.GetUserSettings(guest.ID)
	if _, exists := userSettings["dashboard.units"]; exists {
		t.Error("guest setting should not be persisted")
	}
}

func TestManager_HandleSettingChange_UserCannotPersistGlobalCategory(t *testing.T) {
	store := setupTestStore(t)
	pub := &fakePublisher{}
	mgr := NewManager(store, pub)

	store.SetGlobalSettings(map[string]string{"alarm.enabled": "true"})
	user, _ := store.CreateProfile("User", "user", 1)
	mgr.Activate(user)

	mgr.HandleSettingChange("alarm.enabled", "false")

	global, _ := store.GetGlobalSettings()
	if global["alarm.enabled"] != "true" {
		t.Errorf("alarm.enabled = %q, want %q (user should not change global)", global["alarm.enabled"], "true")
	}

	userSettings, _ := store.GetUserSettings(user.ID)
	if _, exists := userSettings["alarm.enabled"]; exists {
		t.Error("global-category setting should not be persisted as user override")
	}
}
