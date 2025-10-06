package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
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
		p.API.LogError("Uploaded file has zero size")
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

// FileWillBeDownloaded is invoked when a file is about to be downloaded
//
// This demo implementation logs a message when a file is going to be downloaded
// and rejects downloads based on configuration settings.
func (p *Plugin) FileWillBeDownloaded(c *plugin.Context, fileInfo *model.FileInfo, userId string) string {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return ""
	}

	// Log the file download attempt
	p.API.LogInfo("File download attempted",
		"file_name", fileInfo.Name,
		"user_id", userId,
		"file_id", fileInfo.Id)

	// Check if all file downloads should be rejected (testing option)
	if configuration.RejectAllFileDownloads {
		p.API.LogWarn("File download rejected - all downloads disabled for testing",
			"file_name", fileInfo.Name,
			"user_id", userId)

		// Send an ephemeral message to the user who tried to download the file
		rejectionMessage := fmt.Sprintf("Download of file '%s' was rejected. All file downloads are currently disabled for testing purposes.", fileInfo.Name)
		if err := p.sendEphemeralMessage(userId, fileInfo.ChannelId, rejectionMessage); err != nil {
			p.API.LogError("Failed to send download rejection message",
				"user_id", userId,
				"file_name", fileInfo.Name,
				"channel_id", fileInfo.ChannelId,
				"error", err.Error())
		}

		return "All file downloads are currently disabled for testing"
	}

	// Check if the file has a .mp4 extension (case-insensitive)
	if strings.HasSuffix(strings.ToLower(fileInfo.Name), ".mp4") {
		p.API.LogWarn("MP4 file download rejected",
			"file_name", fileInfo.Name,
			"user_id", userId)

		// Send an ephemeral message to the user who tried to download the file
		rejectionMessage := fmt.Sprintf("Download of file '%s' was rejected. MP4 files are not allowed to be downloaded.", fileInfo.Name)
		if err := p.sendEphemeralMessage(userId, fileInfo.ChannelId, rejectionMessage); err != nil {
			p.API.LogError("Failed to send download rejection message",
				"user_id", userId,
				"file_name", fileInfo.Name,
				"channel_id", fileInfo.ChannelId,
				"error", err.Error())
		}

		return "Downloading MP4 files is not allowed"
	}

	// Allow the download for other files
	return ""
}
