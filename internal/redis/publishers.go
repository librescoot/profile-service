package redis

import (
	ipc "github.com/librescoot/redis-ipc"
)

type Publishers struct {
	profile  *ipc.HashPublisher
	settings *ipc.HashPublisher
}

func NewPublishers(client *ipc.Client) *Publishers {
	return &Publishers{
		profile:  client.NewHashPublisher("profile"),
		settings: client.NewHashPublisher("settings"),
	}
}

func (p *Publishers) PublishProfile(fields map[string]string) {
	m := make(map[string]any, len(fields))
	for k, v := range fields {
		m[k] = v
	}
	p.profile.SetMany(m)
}

func (p *Publishers) PublishSettings(fields map[string]string) {
	m := make(map[string]any, len(fields))
	for k, v := range fields {
		m[k] = v
	}
	p.settings.ReplaceAll(m)
}

func (p *Publishers) ClearProfile() {
	p.profile.Clear()
}

func (p *Publishers) SetShowSelector(show bool) {
	val := "false"
	if show {
		val = "true"
	}
	p.profile.Set("show-selector", val)
}
