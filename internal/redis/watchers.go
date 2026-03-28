package redis

import (
	"log/slog"

	ipc "github.com/librescoot/redis-ipc"
)

type ProfileCallbacks struct {
	OnKeycardAuth    func(uid string) error
	OnVehicleState   func(state string) error
	OnSettingChanged func(field, value string) error
}

type Watchers struct {
	keycard  *ipc.HashWatcher
	vehicle  *ipc.HashWatcher
	settings *ipc.HashWatcher
	log      *slog.Logger
}

func NewWatchers(client *ipc.Client, cb *ProfileCallbacks, log *slog.Logger) *Watchers {
	w := &Watchers{log: log}

	w.keycard = client.NewHashWatcher("keycard")
	w.keycard.OnField("authentication", func(value string) error {
		if value != "passed" {
			return nil
		}
		uid, err := w.keycard.Fetch("uid")
		if err != nil {
			log.Error("fetch keycard uid", "error", err)
			return nil
		}
		return cb.OnKeycardAuth(uid)
	})

	w.vehicle = client.NewHashWatcher("vehicle")
	w.vehicle.OnField("state", func(value string) error {
		return cb.OnVehicleState(value)
	})

	w.settings = client.NewHashWatcher("settings")
	w.settings.OnAny(func(field, value string) error {
		return cb.OnSettingChanged(field, value)
	})

	return w
}

func (w *Watchers) Start() error {
	if err := w.keycard.Start(); err != nil {
		return err
	}
	if err := w.vehicle.StartWithSync(); err != nil {
		return err
	}
	if err := w.settings.Start(); err != nil {
		return err
	}
	return nil
}

func (w *Watchers) Stop() {
	w.keycard.Stop()
	w.vehicle.Stop()
	w.settings.Stop()
}
