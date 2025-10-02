// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import (
	"archive/zip"
	"encoding/base32"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769").WithPadding(base32.NoPadding)

// NewID is a globally unique identifier.  It is a [A-Z0-9] string 26
// characters long.  It is a UUID version 4 Guid that is zbased32 encoded
// without the padding.
func NewID() string {
	return encoding.EncodeToString(newRandom()[:])
}

func newRandom() *uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil
	}

	return &id
}

func CoalesceInt(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

func CalculateNPS(promoters, detractors, passives int64) float64 {
	totalResponses := promoters + detractors + passives
	if totalResponses == 0 {
		return 0.0
	}

	nps := (float64(promoters)/float64(totalResponses))*100 - (float64(detractors)/float64(totalResponses))*100
	return nps
}

// CreateZip creates a zip file at the specified `zipFilePath` containing the files listed in `files`.
func CreateZip(zipFilePath string, files []string) error {
	// Create a new zip file
	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Iterate over the files and add them to the zip archive
	for _, file := range files {
		// Open the file to be zipped
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Get the file info
		fileInfo, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		// Create a new zip file header
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		// Set the name of the file within the zip archive
		header.Name = filepath.Base(file)

		// Add the file header to the zip archive
		fileWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the file contents to the zip archive
		_, err = io.Copy(fileWriter, fileToZip)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Zip file created successfully: %s\n", zipFilePath)
	return nil
}
