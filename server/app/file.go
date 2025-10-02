// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (a *WhatsappApp) writeFileLocally(fr io.Reader, path string) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		directory, _ := filepath.Abs(filepath.Dir(path))
		a.api.LogError("writeFileLocally: failed to create dir path", "directory", directory, "path", path, "error", err.Error())
		return 0, errors.Wrapf(err, "writeFileLocally: unable to create the directory %s for the file %s", directory, path)
	}

	fw, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		a.api.LogError("writeFileLocally: unable to open the file to write the data", "filePath", path, "error", err.Error())
		return 0, errors.Wrapf(err, "unable to open the file %s to write the data", path)
	}

	defer fw.Close()

	written, err := io.Copy(fw, fr)
	if err != nil {
		a.api.LogError("writeFileLocally: unable write the data in the file", "filePath", path, "error", err.Error())
		return written, errors.Wrapf(err, "unable write the data in the file %s", path)
	}

	return written, nil
}

func (a *WhatsappApp) readFile(path string) ([]byte, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		a.api.LogError("ReadFile: unable to read file", "path", path, "error", err.Error())
		return nil, errors.Wrapf(err, "ReadFile: unable to read file %s", path)
	}
	return f, nil
}
