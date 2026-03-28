package db

import (
	"fmt"
)

func (s *Store) SetGlobalSettings(settings map[string]string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for key, value := range settings {
		_, err := tx.Exec(
			`INSERT INTO settings_global (key, value) VALUES (?, ?)
			 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
			key, value,
		)
		if err != nil {
			return fmt.Errorf("set global setting %s: %w", key, err)
		}
	}

	return tx.Commit()
}

func (s *Store) SetGlobalSetting(key, value string) error {
	_, err := s.db.Exec(
		`INSERT INTO settings_global (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func (s *Store) GetGlobalSettings() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, value FROM settings_global`)
	if err != nil {
		return nil, fmt.Errorf("get global settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("scan global setting: %w", err)
		}
		settings[key] = value
	}
	return settings, nil
}

func (s *Store) SetUserSettings(profileID string, settings map[string]string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for key, value := range settings {
		_, err := tx.Exec(
			`INSERT INTO settings_user (profile_id, key, value) VALUES (?, ?, ?)
			 ON CONFLICT(profile_id, key) DO UPDATE SET value = excluded.value`,
			profileID, key, value,
		)
		if err != nil {
			return fmt.Errorf("set user setting %s: %w", key, err)
		}
	}

	return tx.Commit()
}

func (s *Store) SetUserSetting(profileID, key, value string) error {
	_, err := s.db.Exec(
		`INSERT INTO settings_user (profile_id, key, value) VALUES (?, ?, ?)
		 ON CONFLICT(profile_id, key) DO UPDATE SET value = excluded.value`,
		profileID, key, value,
	)
	return err
}

func (s *Store) GetUserSettings(profileID string) (map[string]string, error) {
	rows, err := s.db.Query(
		`SELECT key, value FROM settings_user WHERE profile_id = ?`, profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("get user settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("scan user setting: %w", err)
		}
		settings[key] = value
	}
	return settings, nil
}

func (s *Store) DeleteUserSetting(profileID, key string) error {
	_, err := s.db.Exec(
		`DELETE FROM settings_user WHERE profile_id = ? AND key = ?`,
		profileID, key,
	)
	return err
}
