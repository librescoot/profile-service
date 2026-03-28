package profile

import (
	"errors"
	"math/rand/v2"
	"path/filepath"

	"github.com/librescoot/profile-service/internal/avatar"
	"github.com/librescoot/profile-service/internal/db"
	"github.com/librescoot/profile-service/internal/names"
)

var ErrProfileNotFound = errors.New("profile not found")

const DefaultAvatarDir = "/data/profiles/avatars"

type Resolver struct {
	store     *db.Store
	avatarDir string
}

func NewResolver(store *db.Store) *Resolver {
	return &Resolver{store: store, avatarDir: DefaultAvatarDir}
}

func NewResolverWithAvatarDir(store *db.Store, avatarDir string) *Resolver {
	return &Resolver{store: store, avatarDir: avatarDir}
}

func (r *Resolver) ResolveByKeycard(uid string) (*db.Profile, error) {
	p, err := r.store.LookupByIdentifier("keycard", uid)
	if err != nil {
		return nil, ErrProfileNotFound
	}
	return p, nil
}

func (r *Resolver) ResolveByCloudUser(cloudUserID string) (*db.Profile, error) {
	p, err := r.store.LookupByCloudUserID(cloudUserID)
	if err == nil {
		return p, nil
	}

	seed := rand.Int()
	name := names.GenerateWithSeed(seed)
	p, err = r.store.CreateProfile(name, "user", seed)
	if err != nil {
		return nil, err
	}

	if err := r.store.SetCloudUserID(p.ID, cloudUserID); err != nil {
		return nil, err
	}
	p.CloudUserID = cloudUserID

	avatarPath := filepath.Join(r.avatarDir, p.ID+".png")
	avatar.Generate(seed, avatarPath)
	p.AvatarPath = avatarPath

	return p, nil
}

func (r *Resolver) ResolveByBLEMac(mac string) (*db.Profile, error) {
	p, err := r.store.LookupByIdentifier("ble_mac", mac)
	if err != nil {
		return nil, ErrProfileNotFound
	}
	return p, nil
}
