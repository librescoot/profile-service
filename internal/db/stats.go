package db

import (
	"fmt"
	"time"
)

type UsageStats struct {
	ProfileID      string
	TotalDistanceM int64
	TotalTrips     int64
	TotalDurationS int64
	MaxSpeed       float64
	LastRideAt     *int64
}

func (s *Store) GetUsageStats(profileID string) (*UsageStats, error) {
	stats := &UsageStats{ProfileID: profileID}
	var lastRideAt *int64

	err := s.db.QueryRow(
		`SELECT total_distance_m, total_trips, total_duration_s, max_speed, last_ride_at
		 FROM usage_stats WHERE profile_id = ?`, profileID,
	).Scan(&stats.TotalDistanceM, &stats.TotalTrips, &stats.TotalDurationS, &stats.MaxSpeed, &lastRideAt)
	if err != nil {
		return nil, fmt.Errorf("get usage stats: %w", err)
	}

	stats.LastRideAt = lastRideAt
	return stats, nil
}

func (s *Store) AddTripStats(profileID string, distanceM, durationS int64, maxSpeed float64) error {
	now := time.Now().Unix()
	_, err := s.db.Exec(
		`UPDATE usage_stats SET
			total_distance_m = total_distance_m + ?,
			total_trips = total_trips + 1,
			total_duration_s = total_duration_s + ?,
			max_speed = MAX(max_speed, ?),
			last_ride_at = ?
		 WHERE profile_id = ?`,
		distanceM, durationS, maxSpeed, now, profileID,
	)
	return err
}
