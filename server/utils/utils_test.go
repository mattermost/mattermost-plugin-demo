// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewID(t *testing.T) {
	t.Run("generated ID should always be 26 characters long", func(t *testing.T) {
		require.Equal(t, 26, len(NewID()))
	})
}
