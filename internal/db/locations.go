package db

import (
	"fmt"
	"time"
)

type SavedLocation struct {
	ID        int64
	ProfileID string
	Name      string
	Latitude  float64
	Longitude float64
	Type      string
	CreatedAt int64
	UpdatedAt int64
}

func (s *Store) AddSavedLocation(profileID, name string, lat, lon float64, locType string) (*SavedLocation, error) {
	now := time.Now().Unix()
	result, err := s.db.Exec(
		`INSERT INTO saved_locations (profile_id, name, latitude, longitude, type, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		profileID, name, lat, lon, locType, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("add saved location: %w", err)
	}

	id, _ := result.LastInsertId()
	return &SavedLocation{
		ID: id, ProfileID: profileID, Name: name,
		Latitude: lat, Longitude: lon, Type: locType,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (s *Store) ListSavedLocations(profileID string) ([]SavedLocation, error) {
	rows, err := s.db.Query(
		`SELECT id, profile_id, name, latitude, longitude, type, created_at, updated_at
		 FROM saved_locations WHERE profile_id = ? ORDER BY type, name`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("list saved locations: %w", err)
	}
	defer rows.Close()

	var locs []SavedLocation
	for rows.Next() {
		var loc SavedLocation
		if err := rows.Scan(&loc.ID, &loc.ProfileID, &loc.Name, &loc.Latitude, &loc.Longitude, &loc.Type, &loc.CreatedAt, &loc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan location: %w", err)
		}
		locs = append(locs, loc)
	}
	return locs, nil
}

func (s *Store) DeleteSavedLocation(id int64) error {
	_, err := s.db.Exec(`DELETE FROM saved_locations WHERE id = ?`, id)
	return err
}
