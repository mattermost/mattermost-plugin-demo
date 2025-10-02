// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (a *WhatsappApp) CreateSession(userID string) (*model.Session, error) {
	session := &model.Session{
		UserID: userID,
	}

	session.SetDefaults()
	if err := session.IsValid(); err != nil {
		return nil, errors.Wrap(err, "CreateSession: session is not valid")
	}

	if err := a.store.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (a *WhatsappApp) GetSessionByID(sessionID string) (*model.Session, error) {
	session, err := a.store.GetSessionByID(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByID: failed to get session")
	}

	return session, nil
}

func (a *WhatsappApp) UpdateSession(session *model.Session) error {
	if err := session.IsValid(); err != nil {
		return errors.Wrap(err, "UpdateSession: session is not valid")
	}

	if err := a.store.UpdateSession(session); err != nil {
		return errors.Wrap(err, "UpdateSession: failed to update session")
	}

	return nil
}

func (a *WhatsappApp) GetSessions() ([]*model.Session, error) {
	sessions, err := a.store.GetSessions()
	if err != nil {
		return nil, errors.Wrap(err, "GetSessions: failed to get sessions")
	}

	return sessions, nil
}

func (a *WhatsappApp) DeleteSession(sessionID string) error {
	if err := a.store.DeleteSession(sessionID); err != nil {
		return errors.Wrap(err, "DeleteSession: failed to delete session")
	}

	return nil
}

func (a *WhatsappApp) GetActiveSessionByUserID(userID string) (*model.Session, error) {
	session, err := a.store.GetActiveSessionByUserID(userID)
	if err != nil {
		return nil, errors.Wrap(err, "GetActiveSessionByUserID: failed to get active session")
	}
	return session, nil
}
