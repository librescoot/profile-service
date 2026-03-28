package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Identifier struct {
	ID        int
	ProfileID string
	Type      string
	Value     string
	CreatedAt int64
}

func (s *Store) AddIdentifier(profileID, idType, value string) error {
	_, err := s.db.Exec(
		`INSERT INTO identifiers (profile_id, type, value, created_at) VALUES (?, ?, ?, ?)`,
		profileID, idType, value, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("add identifier: %w", err)
	}
	return nil
}

func (s *Store) RemoveIdentifier(profileID, idType, value string) error {
	_, err := s.db.Exec(
		`DELETE FROM identifiers WHERE profile_id = ? AND type = ? AND value = ?`,
		profileID, idType, value,
	)
	return err
}

func (s *Store) LookupByIdentifier(idType, value string) (*Profile, error) {
	var profileID string
	err := s.db.QueryRow(
		`SELECT profile_id FROM identifiers WHERE type = ? AND value = ?`,
		idType, value,
	).Scan(&profileID)
	if err != nil {
		return nil, fmt.Errorf("lookup identifier %s/%s: %w", idType, value, err)
	}
	return s.GetProfile(profileID)
}

func (s *Store) LookupByCloudUserID(cloudUserID string) (*Profile, error) {
	p := &Profile{}
	var cuid, avatarPath sql.NullString
	err := s.db.QueryRow(
		`SELECT id, name, role, cloud_user_id, avatar_type, avatar_seed, avatar_path, created_at, updated_at
		 FROM profiles WHERE cloud_user_id = ?`, cloudUserID,
	).Scan(&p.ID, &p.Name, &p.Role, &cuid, &p.AvatarType, &p.AvatarSeed, &avatarPath, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("lookup cloud user %s: %w", cloudUserID, err)
	}
	p.CloudUserID = cuid.String
	p.AvatarPath = avatarPath.String
	return p, nil
}

func (s *Store) SetCloudUserID(profileID, cloudUserID string) error {
	_, err := s.db.Exec(
		`UPDATE profiles SET cloud_user_id = ?, updated_at = ? WHERE id = ?`,
		cloudUserID, time.Now().Unix(), profileID,
	)
	return err
}

func (s *Store) ListIdentifiers(profileID string) ([]Identifier, error) {
	rows, err := s.db.Query(
		`SELECT id, profile_id, type, value, created_at FROM identifiers WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("list identifiers: %w", err)
	}
	defer rows.Close()

	var ids []Identifier
	for rows.Next() {
		var id Identifier
		if err := rows.Scan(&id.ID, &id.ProfileID, &id.Type, &id.Value, &id.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan identifier: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
