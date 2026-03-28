package db

import (
	"testing"
)

func TestAddIdentifier(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	err := store.AddIdentifier(p.ID, "keycard", "04B1A2C3D4E5F6")
	if err != nil {
		t.Fatalf("AddIdentifier error: %v", err)
	}
}

func TestAddIdentifier_Duplicate(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	store.AddIdentifier(p.ID, "keycard", "04B1A2C3D4E5F6")
	err := store.AddIdentifier(p.ID, "keycard", "04B1A2C3D4E5F6")
	if err == nil {
		t.Error("expected error for duplicate identifier")
	}
}

func TestLookupByIdentifier(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Card User", "user", 1)
	store.AddIdentifier(p.ID, "keycard", "04B1A2C3D4E5F6")

	got, err := store.LookupByIdentifier("keycard", "04B1A2C3D4E5F6")
	if err != nil {
		t.Fatalf("LookupByIdentifier error: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("got ID %q, want %q", got.ID, p.ID)
	}
}

func TestLookupByIdentifier_NotFound(t *testing.T) {
	store := openTestDB(t)

	_, err := store.LookupByIdentifier("keycard", "DEADBEEF")
	if err == nil {
		t.Error("expected error for unknown identifier")
	}
}

func TestLookupByCloudUserID(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Cloud User", "user", 1)
	store.SetCloudUserID(p.ID, "cloud-123")

	got, err := store.LookupByCloudUserID("cloud-123")
	if err != nil {
		t.Fatalf("LookupByCloudUserID error: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("got ID %q, want %q", got.ID, p.ID)
	}
}

func TestRemoveIdentifier(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)
	store.AddIdentifier(p.ID, "keycard", "AABBCCDD")

	err := store.RemoveIdentifier(p.ID, "keycard", "AABBCCDD")
	if err != nil {
		t.Fatalf("RemoveIdentifier error: %v", err)
	}

	_, err = store.LookupByIdentifier("keycard", "AABBCCDD")
	if err == nil {
		t.Error("expected error after removal")
	}
}

func TestListIdentifiers(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)
	store.AddIdentifier(p.ID, "keycard", "CARD1")
	store.AddIdentifier(p.ID, "ble_mac", "AA:BB:CC:DD:EE:FF")

	ids, err := store.ListIdentifiers(p.ID)
	if err != nil {
		t.Fatalf("ListIdentifiers error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("len = %d, want 2", len(ids))
	}
}
