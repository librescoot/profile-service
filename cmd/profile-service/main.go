package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/librescoot/profile-service/internal/app"
)

var version = "dev"

func main() {
	redisAddr := flag.String("redis", "192.168.7.1:6379", "Redis address")
	dbPath := flag.String("db", "/data/profiles.db", "SQLite database path")
	tomlPath := flag.String("toml", "/data/settings.toml", "Path to legacy settings TOML for migration")
	avatarDir := flag.String("avatar-dir", "/data/profiles/avatars", "Avatar image directory")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("profile-service %s\n", version)
		os.Exit(0)
	}

	level := slog.LevelInfo
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}
	if os.Getenv("JOURNAL_STREAM") != "" {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	application := app.New(&app.Config{
		RedisAddr: *redisAddr,
		DBPath:    *dbPath,
		TOMLPath:  *tomlPath,
		AvatarDir: *avatarDir,
		Logger:    logger,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- application.Run(ctx)
	}()

	select {
	case sig := <-sigChan:
		logger.Info("received signal", "signal", sig)
		cancel()
		<-errChan
	case err := <-errChan:
		if err != nil {
			logger.Error("application error", "error", err)
			os.Exit(1)
		}
	}
}
