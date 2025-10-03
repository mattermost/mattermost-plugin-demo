// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (app *WhatsappApp) CreateSession(userID string) (*model.Session, error) {
	session := &model.Session{
		UserID: userID,
	}

	session.SetDefaults()
	if err := session.IsValid(); err != nil {
		return nil, errors.Wrap(err, "CreateSession: session is not valid")
	}

	if err := app.store.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (app *WhatsappApp) GetSessionByID(sessionID string) (*model.Session, error) {
	session, err := app.store.GetSessionByID(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByID: failed to get session")
	}

	return session, nil
}

func (app *WhatsappApp) UpdateSession(session *model.Session) error {
	if err := session.IsValid(); err != nil {
		return errors.Wrap(err, "UpdateSession: session is not valid")
	}

	if err := app.store.UpdateSession(session); err != nil {
		return errors.Wrap(err, "UpdateSession: failed to update session")
	}

	return nil
}

func (app *WhatsappApp) GetSessions() ([]*model.Session, error) {
	sessions, err := app.store.GetSessions()
	if err != nil {
		return nil, errors.Wrap(err, "GetSessions: failed to get sessions")
	}

	return sessions, nil
}

func (app *WhatsappApp) DeleteSession(sessionID string) error {
	if err := app.store.DeleteSession(sessionID); err != nil {
		return errors.Wrap(err, "DeleteSession: failed to delete session")
	}

	return nil
}
