package main

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func (p *Plugin) GenerateSupportData(_ *plugin.Context) ([]*model.FileData, error) {
	p.API.LogInfo("Generating demo support data")
	b := []byte("this is a demo file.")

	return []*model.FileData{{
		Filename: "demo_plugin_support_packet.txt",
		Body:     b,
	},
	}, nil
}
