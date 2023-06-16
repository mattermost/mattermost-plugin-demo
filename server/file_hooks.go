package main

import (
	"bytes"
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

// FileWillBeUploaded is invoked when a file is uploaded, but before it is committed to backing store
//
// This demo implementation logs a message to the demo channel in the team
// when a new file is uploaded.
func (p *Plugin) FileWillBeUploaded(c *plugin.Context, fileInfo *model.FileInfo, reader bytes.Reader, buf *bytes.Buffer) (*model.FileInfo, string) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return nil, "Configuration is disabled"
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"Failed to query teams FileWillBeUploaded",
			"error", err.Error(),
		)
		return nil, "Failed to query teams"
	}

	if reader.Size() == 0 {
		p.API.LogError(
			"Uploaded file has zero size",
			"error", err.Error(),
		)
		return nil, "Upload Failed as file has zero size"
	}

	for _, team := range teams {
		msg := fmt.Sprintf("FileName @%s has been created in", fileInfo.Name)
		if err := p.postPluginMessage(team.Id, msg); err != nil {
			p.API.LogError(
				"Failed to post FileWillBeUploaded message",
				"channel_id", configuration.demoChannelIDs[team.Id],
				"error", err.Error(),
			)
		}
	}
	return nil, ""
}
