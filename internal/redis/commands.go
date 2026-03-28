package redis

import (
	"fmt"
	"log/slog"
	"strings"

	ipc "github.com/librescoot/redis-ipc"
)

type CommandCallbacks struct {
	OnSelectProfile func(id string) error
	OnRefresh       func() error
}

type CommandHandler struct {
	handler *ipc.QueueHandler[string]
	log     *slog.Logger
}

func NewCommandHandler(client *ipc.Client, cb *CommandCallbacks, log *slog.Logger) *CommandHandler {
	h := &CommandHandler{log: log}

	h.handler = ipc.HandleRequests(client, "scooter:profile", func(cmd string) error {
		log.Info("received command", "cmd", cmd)

		switch {
		case strings.HasPrefix(cmd, "select:"):
			id := strings.TrimPrefix(cmd, "select:")
			return cb.OnSelectProfile(id)
		case cmd == "refresh":
			return cb.OnRefresh()
		default:
			return fmt.Errorf("unknown profile command: %s", cmd)
		}
	})

	return h
}

func (c *CommandHandler) Stop() {
	if c.handler != nil {
		c.handler.Stop()
	}
}
