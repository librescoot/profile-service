# profile-service

Figures out who's riding the scooter and applies their settings. Replaces `settings-service`.

Part of the [Librescoot](https://librescoot.org/) open-source platform.

```bash
GOTOOLCHAIN=go1.25.7 make build-host
./bin/profile-service --redis localhost:6379 --db /tmp/profiles.db
```

## What it does

- Identifies the current rider via keycard UID, cloud user ID, or BLE MAC
- Manages per-user profiles with settings, saved locations, and usage stats
- Merges per-user settings on top of global defaults into the Redis `settings` hash (transparent to all other services)
- Enforces access roles: owner can change anything, shared users get their own overrides for dashboard/nav settings, guests get session-only changes that aren't persisted
- Auto-generates profiles with cute names (Adjective + Animal) and parametric avatars

## How identification works

On unlock, the service checks (in order):

1. Keycard tap? Look up the UID in the identifiers table
2. Cloud unlock? Match or auto-create a profile from the cloud user ID
3. BLE connection? Look up the phone's MAC (requires nRF52 firmware support, not yet implemented)
4. No match? Show a profile selector on the dashboard, or fall back to guest mode

On lock, per-user settings are reverted to global-only.

## Storage

SQLite WAL at `/data/profiles.db`. Six tables: `profiles`, `identifiers`, `settings_global`, `settings_user`, `saved_locations`, `usage_stats`.

Survives reboots and OTA updates (`/data/` partition).

## Redis interface

| Key | Type | What's in it |
|-----|------|-------------|
| `profile` | Hash | Active rider: `active.id`, `active.name`, `active.role`, `active.avatar-path`, `show-selector` |
| `settings` | Hash | Merged settings (global + active user overrides). Other services read this as before. |
| `scooter:profile` | List | Command queue: `select:<id>`, `refresh` |

Watches `keycard` (auth events), `vehicle` (state transitions), `settings` (change tracking for persistence).

## Settings categories

Global-only (owner persists, others session-only): `scooter.*`, `alarm.*`, `engine-ecu.*`, `keycard.*`, `cellular.*`, `updates.*`

Per-user overridable (persisted for shared users): `dashboard.*`, `navigation.*`

## First boot

Creates an owner profile with a generated name and avatar. If `/data/settings.toml` exists (from the old settings-service), imports it into SQLite.

## Build

```bash
GOTOOLCHAIN=go1.25.7 make build        # ARM cross-compile
GOTOOLCHAIN=go1.25.7 make build-host   # native
GOTOOLCHAIN=go1.25.7 make test
```

## Deploy

```bash
scp bin/profile-service deep-blue:/data/profile-service-test
ssh deep-blue "systemctl stop librescoot-profile && cp /data/profile-service-test /usr/bin/profile-service && systemctl start librescoot-profile"
```

## License

This project is dual-licensed. The source code is available under the
[Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License][cc-by-nc-sa].
The maintainers reserve the right to grant separate licenses for commercial distribution; please contact the maintainers to discuss commercial licensing.

[![CC BY-NC-SA 4.0][cc-by-nc-sa-image]][cc-by-nc-sa]

[cc-by-nc-sa]: http://creativecommons.org/licenses/by-nc-sa/4.0/
[cc-by-nc-sa-image]: https://licensebuttons.net/l/by-nc-sa/4.0/88x31.png
