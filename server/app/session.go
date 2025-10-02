package app

import (
	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	"github.com/pkg/errors"
)

func (a *WhatsappApp) GetSession(sessionId string) (*model.WhatsappSession, error) {
	session, err := a.store.GetSession(sessionId)
	if err != nil {
		return nil, errors.Wrap(err, "GetSession: failed to get session from database")
	}
	return session, nil
}

func (a *WhatsappApp) CreateSession(userId string) (*model.WhatsappSession, error) {
	session, err := a.store.CreateSession(userId)
	if err != nil {
		return nil, errors.Wrap(err, "CreateSession: failed to create session")
	}
	return session, nil
}

func (a *WhatsappApp) CloseSession(sessionId string) (*model.WhatsappSession, error) {
	session, err := a.store.CloseSession(sessionId)
	if err != nil {
		return nil, errors.Wrap(err, "CloseSession: failed to close session")
	}
	return session, nil
}

func (a *WhatsappApp) GetSessionByUserId(userId string) (*model.WhatsappSession, error) {
	session, err := a.store.GetSessionByUserId(userId)
	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByUserId: failed to get session by user id")
	}
	return session, nil
}

func (a *WhatsappApp) GetSessionsUnclosed() ([]model.WhatsappSession, error) {
	sessions, err := a.store.GetSessionsUnclosed()
	if err != nil {
		return nil, errors.Wrap(err, "GetSessionsUnclosed: failed to get unclosed sessions")
	}
	return sessions, nil
}
