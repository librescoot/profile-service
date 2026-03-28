package profile

import (
	"github.com/librescoot/profile-service/internal/db"
	"github.com/librescoot/profile-service/internal/settings"
)

type Publisher interface {
	PublishProfile(fields map[string]string)
	PublishSettings(fields map[string]string)
	ClearProfile()
}

type Manager struct {
	store  *db.Store
	pub    Publisher
	active *db.Profile
}

func NewManager(store *db.Store, pub Publisher) *Manager {
	return &Manager{
		store: store,
		pub:   pub,
	}
}

func (m *Manager) ActiveProfile() *db.Profile {
	return m.active
}

func (m *Manager) Activate(p *db.Profile) error {
	m.active = p

	m.pub.PublishProfile(map[string]string{
		"active.id":            p.ID,
		"active.name":          p.Name,
		"active.role":          p.Role,
		"active.avatar-path":   p.AvatarPath,
		"active.cloud-user-id": p.CloudUserID,
	})

	return m.publishMergedSettings()
}

func (m *Manager) Deactivate() error {
	m.active = nil
	m.pub.ClearProfile()
	return m.publishGlobalSettings()
}

func (m *Manager) HandleSettingChange(key, value string) {
	if m.active == nil {
		return
	}

	switch m.active.Role {
	case "owner":
		m.store.SetGlobalSetting(key, value)
	case "user":
		if settings.IsUserOverridable(key) {
			m.store.SetUserSetting(m.active.ID, key, value)
		}
	case "guest":
	}
}

func (m *Manager) publishMergedSettings() error {
	global, err := m.store.GetGlobalSettings()
	if err != nil {
		return err
	}

	var userOverrides map[string]string
	if m.active != nil {
		userOverrides, err = m.store.GetUserSettings(m.active.ID)
		if err != nil {
			return err
		}
	}

	merged := settings.Merge(global, userOverrides)
	m.pub.PublishSettings(merged)
	return nil
}

func (m *Manager) publishGlobalSettings() error {
	global, err := m.store.GetGlobalSettings()
	if err != nil {
		return err
	}

	m.pub.PublishSettings(global)
	return nil
}
