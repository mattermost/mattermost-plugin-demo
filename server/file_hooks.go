package main

import (
	"bytes"
	"fmt"

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
// The downloadType parameter indicates what type of access is being attempted and can be:
//   - model.FileDownloadTypeFile: Full file download
//   - model.FileDownloadTypeThumbnail: Thumbnail request
//   - model.FileDownloadTypePreview: Preview image request
//   - model.FileDownloadTypePublic: Public link access (userId will be empty in this case)
func (p *Plugin) FileWillBeDownloaded(c *plugin.Context, fileInfo *model.FileInfo, userId string, downloadType model.FileDownloadType) string {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return ""
	}

	// Log the file download attempt
	p.API.LogInfo("File download attempted",
		"file_name", fileInfo.Name,
		"user_id", userId,
		"download_type", string(downloadType),
		"file_id", fileInfo.Id)

	// Check if all file downloads should be rejected (testing option)
	if configuration.RejectAllFileDownloads {
		return "All file downloads are currently disabled by the demo plugin"
	}

	// Example: Apply different policies based on download type
	// Uncomment and modify as needed:
	// if downloadType == model.FileDownloadTypePublic {
	// 	return "Public file downloads are disabled"
	// }
	// if downloadType == model.FileDownloadTypePreview && someCondition {
	// 	return "Preview downloads are restricted"
	// }

	// Allow the download for other files
	return ""
}
