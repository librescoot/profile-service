package db

import (
	"testing"
)

func TestGetUsageStats(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	stats, err := store.GetUsageStats(p.ID)
	if err != nil {
		t.Fatalf("GetUsageStats error: %v", err)
	}
	if stats.TotalTrips != 0 {
		t.Errorf("total_trips = %d, want 0", stats.TotalTrips)
	}
}

func TestUpdateUsageStats(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	err := store.AddTripStats(p.ID, 5000, 300, 45.5)
	if err != nil {
		t.Fatalf("AddTripStats error: %v", err)
	}

	stats, _ := store.GetUsageStats(p.ID)
	if stats.TotalDistanceM != 5000 {
		t.Errorf("distance = %d, want 5000", stats.TotalDistanceM)
	}
	if stats.TotalTrips != 1 {
		t.Errorf("trips = %d, want 1", stats.TotalTrips)
	}
	if stats.TotalDurationS != 300 {
		t.Errorf("duration = %d, want 300", stats.TotalDurationS)
	}
	if stats.MaxSpeed != 45.5 {
		t.Errorf("max_speed = %f, want 45.5", stats.MaxSpeed)
	}
}

func TestAddTripStats_AccumulatesMultipleTrips(t *testing.T) {
	store := openTestDB(t)
	p, _ := store.CreateProfile("Test", "user", 1)

	store.AddTripStats(p.ID, 1000, 120, 30.0)
	store.AddTripStats(p.ID, 2000, 180, 50.0)

	stats, _ := store.GetUsageStats(p.ID)
	if stats.TotalDistanceM != 3000 {
		t.Errorf("distance = %d, want 3000", stats.TotalDistanceM)
	}
	if stats.TotalTrips != 2 {
		t.Errorf("trips = %d, want 2", stats.TotalTrips)
	}
	if stats.MaxSpeed != 50.0 {
		t.Errorf("max_speed = %f, want 50.0 (should keep max)", stats.MaxSpeed)
	}
}
