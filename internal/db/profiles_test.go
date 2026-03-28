package db

import (
	"path/filepath"
	"testing"
)

func openTestDB(t *testing.T) *Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestCreateProfile(t *testing.T) {
	store := openTestDB(t)

	p, err := store.CreateProfile("Swift Otter", "owner", 42)
	if err != nil {
		t.Fatalf("CreateProfile error: %v", err)
	}
	if p.ID == "" {
		t.Error("expected non-empty ID")
	}
	if p.Name != "Swift Otter" {
		t.Errorf("name = %q, want %q", p.Name, "Swift Otter")
	}
	if p.Role != "owner" {
		t.Errorf("role = %q, want %q", p.Role, "owner")
	}
	if p.AvatarSeed != 42 {
		t.Errorf("avatar_seed = %d, want %d", p.AvatarSeed, 42)
	}
}

func TestGetProfile(t *testing.T) {
	store := openTestDB(t)

	created, err := store.CreateProfile("Calm Falcon", "user", 99)
	if err != nil {
		t.Fatalf("CreateProfile error: %v", err)
	}

	got, err := store.GetProfile(created.ID)
	if err != nil {
		t.Fatalf("GetProfile error: %v", err)
	}
	if got.Name != "Calm Falcon" {
		t.Errorf("name = %q, want %q", got.Name, "Calm Falcon")
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	store := openTestDB(t)

	_, err := store.GetProfile("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent profile")
	}
}

func TestListProfiles(t *testing.T) {
	store := openTestDB(t)

	store.CreateProfile("Alpha", "owner", 1)
	store.CreateProfile("Beta", "user", 2)

	profiles, err := store.ListProfiles()
	if err != nil {
		t.Fatalf("ListProfiles error: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("len = %d, want 2", len(profiles))
	}
}

func TestUpdateProfile(t *testing.T) {
	store := openTestDB(t)

	p, _ := store.CreateProfile("Old Name", "user", 1)
	err := store.UpdateProfile(p.ID, "New Name", "")
	if err != nil {
		t.Fatalf("UpdateProfile error: %v", err)
	}

	got, _ := store.GetProfile(p.ID)
	if got.Name != "New Name" {
		t.Errorf("name = %q, want %q", got.Name, "New Name")
	}
}

func TestDeleteProfile(t *testing.T) {
	store := openTestDB(t)

	p, _ := store.CreateProfile("Doomed", "guest", 1)
	err := store.DeleteProfile(p.ID)
	if err != nil {
		t.Fatalf("DeleteProfile error: %v", err)
	}

	_, err = store.GetProfile(p.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestGetOwnerProfile(t *testing.T) {
	store := openTestDB(t)

	store.CreateProfile("Not Owner", "user", 1)
	owner, _ := store.CreateProfile("The Owner", "owner", 2)

	got, err := store.GetOwnerProfile()
	if err != nil {
		t.Fatalf("GetOwnerProfile error: %v", err)
	}
	if got.ID != owner.ID {
		t.Errorf("got ID %q, want %q", got.ID, owner.ID)
	}
}
