package app

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"

	"github.com/librescoot/profile-service/internal/avatar"
	"github.com/librescoot/profile-service/internal/db"
	"github.com/librescoot/profile-service/internal/names"
	"github.com/librescoot/profile-service/internal/profile"
	profileRedis "github.com/librescoot/profile-service/internal/redis"
	"github.com/librescoot/profile-service/internal/settings"
	ipc "github.com/librescoot/redis-ipc"
)

type Config struct {
	RedisAddr string
	DBPath    string
	TOMLPath  string
	AvatarDir string
	Logger    *slog.Logger
}

type App struct {
	cfg      *Config
	log      *slog.Logger
	store    *db.Store
	client   *ipc.Client
	pubs     *profileRedis.Publishers
	watchers *profileRedis.Watchers
	cmdH     *profileRedis.CommandHandler
	resolver *profile.Resolver
	manager  *profile.Manager
}

func New(cfg *Config) *App {
	return &App{
		cfg: cfg,
		log: cfg.Logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting profile-service", "redis", a.cfg.RedisAddr, "db", a.cfg.DBPath)

	store, err := db.Open(a.cfg.DBPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.Close()
	a.store = store

	if err := a.ensureOwnerProfile(); err != nil {
		return fmt.Errorf("ensure owner: %w", err)
	}

	if err := a.migrateSettingsIfNeeded(); err != nil {
		a.log.Warn("settings migration failed", "error", err)
	}

	client, err := ipc.New(
		ipc.WithURL(a.cfg.RedisAddr),
		ipc.WithCodec(ipc.StringCodec{}),
		ipc.WithLogger(a.log),
	)
	if err != nil {
		return fmt.Errorf("create redis client: %w", err)
	}
	defer client.Close()
	a.client = client

	a.pubs = profileRedis.NewPublishers(client)
	a.resolver = profile.NewResolverWithAvatarDir(store, a.cfg.AvatarDir)
	a.manager = profile.NewManager(store, a.pubs)

	if err := a.publishGlobalSettings(); err != nil {
		return fmt.Errorf("publish initial settings: %w", err)
	}

	a.watchers = profileRedis.NewWatchers(client, &profileRedis.ProfileCallbacks{
		OnKeycardAuth:    a.handleKeycardAuth,
		OnVehicleState:   a.handleVehicleState,
		OnSettingChanged: a.handleSettingChanged,
	}, a.log)

	a.cmdH = profileRedis.NewCommandHandler(client, &profileRedis.CommandCallbacks{
		OnSelectProfile: a.handleSelectProfile,
		OnRefresh:       a.handleRefresh,
	}, a.log)

	if err := a.watchers.Start(); err != nil {
		return fmt.Errorf("start watchers: %w", err)
	}

	a.log.Info("profile-service ready")
	<-ctx.Done()

	a.log.Info("shutting down")
	a.watchers.Stop()
	a.cmdH.Stop()
	return nil
}

func (a *App) ensureOwnerProfile() error {
	_, err := a.store.GetOwnerProfile()
	if err == nil {
		return nil
	}

	seed := rand.Int()
	name := names.GenerateWithSeed(seed)
	p, err := a.store.CreateProfile(name, "owner", seed)
	if err != nil {
		return err
	}

	avatarPath := a.cfg.AvatarDir + "/" + p.ID + ".png"
	if err := avatar.Generate(seed, avatarPath); err != nil {
		a.log.Warn("avatar generation failed", "error", err)
	}

	a.log.Info("created owner profile", "id", p.ID, "name", name)
	return nil
}

func (a *App) migrateSettingsIfNeeded() error {
	existing, err := a.store.GetGlobalSettings()
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	if _, err := os.Stat(a.cfg.TOMLPath); os.IsNotExist(err) {
		return nil
	}

	imported, err := settings.ImportTOML(a.cfg.TOMLPath)
	if err != nil {
		return err
	}

	if err := a.store.SetGlobalSettings(imported); err != nil {
		return err
	}

	a.log.Info("migrated settings from TOML", "count", len(imported))
	return nil
}

func (a *App) publishGlobalSettings() error {
	global, err := a.store.GetGlobalSettings()
	if err != nil {
		return err
	}
	a.pubs.PublishSettings(global)
	return nil
}

func (a *App) handleKeycardAuth(uid string) error {
	a.log.Info("keycard auth", "uid", uid)

	p, err := a.resolver.ResolveByKeycard(uid)
	if err != nil {
		a.log.Info("unmapped keycard, entering guest mode", "uid", uid)
		return nil
	}

	return a.manager.Activate(p)
}

func (a *App) handleVehicleState(state string) error {
	a.log.Info("vehicle state", "state", state)

	switch state {
	case "parked", "stand-by", "hibernation", "shutting-down":
		if a.manager.ActiveProfile() != nil {
			return a.manager.Deactivate()
		}
	}
	return nil
}

func (a *App) handleSettingChanged(field, value string) error {
	a.manager.HandleSettingChange(field, value)
	return nil
}

func (a *App) handleSelectProfile(id string) error {
	a.log.Info("profile selected", "id", id)

	p, err := a.store.GetProfile(id)
	if err != nil {
		return err
	}

	a.pubs.SetShowSelector(false)
	return a.manager.Activate(p)
}

func (a *App) handleRefresh() error {
	if a.manager.ActiveProfile() != nil {
		return a.manager.Activate(a.manager.ActiveProfile())
	}
	return a.publishGlobalSettings()
}
