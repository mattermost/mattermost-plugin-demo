// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/utils"
)

type Session struct {
	ID       string `json:"id"`
	UserID   string `json:"userID"`
	CreateAt int64  `json:"createAt"`
	ClosedAt *int64 `json:"closedAt,omitempty"`
}

func (s *Session) SetDefaults() {
	if s.ID == "" {
		s.ID = utils.NewID()
	}

	if s.CreateAt == 0 {
		s.CreateAt = mmModel.GetMillis()
	}
}

func (s *Session) IsValid() error {
	if s.ID == "" {
		return errors.New("session ID cannot be empty")
	}

	if s.UserID == "" {
		return errors.New("session user ID cannot be empty")
	}

	if s.CreateAt == 0 {
		return errors.New("session create at time cannot be empty")
	}

	return nil
}
