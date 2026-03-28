package profile

import (
	"path/filepath"
	"testing"

	"github.com/librescoot/profile-service/internal/db"
)

func setupTestStore(t *testing.T) *db.Store {
	t.Helper()
	store, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestResolveByKeycard_Found(t *testing.T) {
	store := setupTestStore(t)
	p, _ := store.CreateProfile("Card User", "user", 1)
	store.AddIdentifier(p.ID, "keycard", "04AABBCCDD")

	resolver := NewResolver(store)
	got, err := resolver.ResolveByKeycard("04AABBCCDD")
	if err != nil {
		t.Fatalf("ResolveByKeycard error: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("got ID %q, want %q", got.ID, p.ID)
	}
}

func TestResolveByKeycard_NotFound(t *testing.T) {
	store := setupTestStore(t)
	resolver := NewResolver(store)

	_, err := resolver.ResolveByKeycard("UNKNOWN")
	if err != ErrProfileNotFound {
		t.Errorf("err = %v, want ErrProfileNotFound", err)
	}
}

func TestResolveByCloudUser_Found(t *testing.T) {
	store := setupTestStore(t)
	p, _ := store.CreateProfile("Cloud User", "user", 1)
	store.SetCloudUserID(p.ID, "cloud-abc")

	resolver := NewResolver(store)
	got, err := resolver.ResolveByCloudUser("cloud-abc")
	if err != nil {
		t.Fatalf("ResolveByCloudUser error: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("got ID %q, want %q", got.ID, p.ID)
	}
}

func TestResolveByCloudUser_NotFound_CreatesProfile(t *testing.T) {
	store := setupTestStore(t)
	resolver := NewResolverWithAvatarDir(store, t.TempDir())

	got, err := resolver.ResolveByCloudUser("new-cloud-user")
	if err != nil {
		t.Fatalf("ResolveByCloudUser error: %v", err)
	}
	if got.CloudUserID != "new-cloud-user" {
		t.Errorf("cloud_user_id = %q, want %q", got.CloudUserID, "new-cloud-user")
	}
	if got.Role != "user" {
		t.Errorf("role = %q, want %q", got.Role, "user")
	}
}

func TestResolveByBLEMac_Found(t *testing.T) {
	store := setupTestStore(t)
	p, _ := store.CreateProfile("BLE User", "user", 1)
	store.AddIdentifier(p.ID, "ble_mac", "AA:BB:CC:DD:EE:FF")

	resolver := NewResolver(store)
	got, err := resolver.ResolveByBLEMac("AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("ResolveByBLEMac error: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("got ID %q, want %q", got.ID, p.ID)
	}
}
