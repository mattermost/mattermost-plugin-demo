// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var UnsafeContentTypes = [...]string{
	"application/javascript",
	"application/ecmascript",
	"text/javascript",
	"text/ecmascript",
	"application/x-javascript",
	"text/html",
}

var MediaContentTypes = [...]string{
	"image/jpeg",
	"image/png",
	"image/bmp",
	"image/gif",
	"image/tiff",
	"image/webp",
	"video/avi",
	"video/mpeg",
	"video/mp4",
	"audio/mpeg",
	"audio/wav",
}

// WriteFileResponse is copied from Mattermost server https://github.com/mattermost/mattermost/blob/f0121d4f23e065a05e63fd452486570271c6a28f/server/platform/shared/web/files.go#L38
func WriteFileResponse(filename string, contentType string, contentSize int64, lastModification time.Time, webserverMode string, fileReader io.ReadSeeker, forceDownload bool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if contentSize > 0 {
		contentSizeStr := strconv.Itoa(int(contentSize))
		if webserverMode == "gzip" {
			w.Header().Set("X-Uncompressed-Content-Length", contentSizeStr)
		} else {
			w.Header().Set("Content-Length", contentSizeStr)
		}
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	} else {
		for _, unsafeContentType := range UnsafeContentTypes {
			if strings.HasPrefix(contentType, unsafeContentType) {
				contentType = "text/plain"
				break
			}
		}
	}

	w.Header().Set("Content-Type", contentType)

	var toDownload bool
	if forceDownload {
		toDownload = true
	} else {
		isMediaType := false

		for _, mediaContentType := range MediaContentTypes {
			if strings.HasPrefix(contentType, mediaContentType) {
				isMediaType = true
				break
			}
		}

		toDownload = !isMediaType
	}

	filename = url.PathEscape(filename)

	if toDownload {
		w.Header().Set("Content-Disposition", "attachment;filename=\""+filename+"\"; filename*=UTF-8''"+filename)
	} else {
		w.Header().Set("Content-Disposition", "inline;filename=\""+filename+"\"; filename*=UTF-8''"+filename)
	}

	// prevent file links from being embedded in iframes
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "Frame-ancestors 'none'")

	http.ServeContent(w, r, filename, lastModification, fileReader)
}
