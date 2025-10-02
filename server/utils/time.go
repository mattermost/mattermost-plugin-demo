// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import "time"

func FormatUnixTimeMillis(timestamp int64) string {
	// Convert milliseconds to seconds
	seconds := timestamp / 1000

	// Convert Unix timestamp to time.Time
	t := time.Unix(seconds, 0)

	// Format the time as "Day Month Year"
	formattedDate := t.Format("2 January 2006")

	return formattedDate
}
