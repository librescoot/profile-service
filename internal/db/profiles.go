package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

type Profile struct {
	ID          string
	Name        string
	Role        string
	CloudUserID string
	AvatarType  string
	AvatarSeed  int
	AvatarPath  string
	CreatedAt   int64
	UpdatedAt   int64
}

func newID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Store) CreateProfile(name, role string, avatarSeed int) (*Profile, error) {
	p := &Profile{
		ID:         newID(),
		Name:       name,
		Role:       role,
		AvatarType: "parametric",
		AvatarSeed: avatarSeed,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	_, err := s.db.Exec(
		`INSERT INTO profiles (id, name, role, avatar_type, avatar_seed, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Role, p.AvatarType, p.AvatarSeed, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert profile: %w", err)
	}

	_, err = s.db.Exec(
		`INSERT INTO usage_stats (profile_id) VALUES (?)`, p.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert usage_stats: %w", err)
	}

	return p, nil
}

func (s *Store) GetProfile(id string) (*Profile, error) {
	p := &Profile{}
	var cloudUserID, avatarPath sql.NullString

	err := s.db.QueryRow(
		`SELECT id, name, role, cloud_user_id, avatar_type, avatar_seed, avatar_path, created_at, updated_at
		 FROM profiles WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Role, &cloudUserID, &p.AvatarType, &p.AvatarSeed, &avatarPath, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get profile %s: %w", id, err)
	}

	p.CloudUserID = cloudUserID.String
	p.AvatarPath = avatarPath.String
	return p, nil
}

func (s *Store) GetOwnerProfile() (*Profile, error) {
	p := &Profile{}
	var cloudUserID, avatarPath sql.NullString

	err := s.db.QueryRow(
		`SELECT id, name, role, cloud_user_id, avatar_type, avatar_seed, avatar_path, created_at, updated_at
		 FROM profiles WHERE role = 'owner' LIMIT 1`,
	).Scan(&p.ID, &p.Name, &p.Role, &cloudUserID, &p.AvatarType, &p.AvatarSeed, &avatarPath, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get owner profile: %w", err)
	}

	p.CloudUserID = cloudUserID.String
	p.AvatarPath = avatarPath.String
	return p, nil
}

func (s *Store) ListProfiles() ([]Profile, error) {
	rows, err := s.db.Query(
		`SELECT id, name, role, cloud_user_id, avatar_type, avatar_seed, avatar_path, created_at, updated_at
		 FROM profiles ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		var cloudUserID, avatarPath sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Role, &cloudUserID, &p.AvatarType, &p.AvatarSeed, &avatarPath, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan profile: %w", err)
		}
		p.CloudUserID = cloudUserID.String
		p.AvatarPath = avatarPath.String
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (s *Store) UpdateProfile(id, name, avatarPath string) error {
	now := time.Now().Unix()
	if avatarPath != "" {
		_, err := s.db.Exec(
			`UPDATE profiles SET name = ?, avatar_type = 'custom', avatar_path = ?, updated_at = ? WHERE id = ?`,
			name, avatarPath, now, id,
		)
		return err
	}
	_, err := s.db.Exec(
		`UPDATE profiles SET name = ?, updated_at = ? WHERE id = ?`,
		name, now, id,
	)
	return err
}

func (s *Store) DeleteProfile(id string) error {
	result, err := s.db.Exec(`DELETE FROM profiles WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("profile %s not found", id)
	}
	return nil
}
