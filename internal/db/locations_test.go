package db

import (
	"testing"
)

func TestAddSavedLocation(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	loc, err := store.AddSavedLocation(p.ID, "Home", 52.5200, 13.4050, "home")
	if err != nil {
		t.Fatalf("AddSavedLocation error: %v", err)
	}
	if loc.Name != "Home" {
		t.Errorf("name = %q, want %q", loc.Name, "Home")
	}
}

func TestListSavedLocations(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	store.AddSavedLocation(p.ID, "Home", 52.52, 13.40, "home")
	store.AddSavedLocation(p.ID, "Work", 52.51, 13.39, "work")

	locs, err := store.ListSavedLocations(p.ID)
	if err != nil {
		t.Fatalf("ListSavedLocations error: %v", err)
	}
	if len(locs) != 2 {
		t.Errorf("len = %d, want 2", len(locs))
	}
}

func TestDeleteSavedLocation(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	loc, _ := store.AddSavedLocation(p.ID, "Temp", 0, 0, "favorite")
	err := store.DeleteSavedLocation(loc.ID)
	if err != nil {
		t.Fatalf("DeleteSavedLocation error: %v", err)
	}

	locs, _ := store.ListSavedLocations(p.ID)
	if len(locs) != 0 {
		t.Errorf("len = %d, want 0", len(locs))
	}
}

func TestSavedLocations_DeletedWithProfile(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)
	store.AddSavedLocation(p.ID, "Home", 52.52, 13.40, "home")

	store.DeleteProfile(p.ID)

	locs, _ := store.ListSavedLocations(p.ID)
	if len(locs) != 0 {
		t.Errorf("len = %d, want 0 after profile delete", len(locs))
	}
}
