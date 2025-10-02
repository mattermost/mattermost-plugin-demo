package app

import (
	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	"github.com/pkg/errors"
)

func (a *WhatsappApp) GetWhatsappChannels() ([]model.WhatsappChannel, error) {
	channels, err := a.store.GetWhatsappChannels()
	if err != nil {
		return nil, errors.Wrap(err, "GetWhatsappChannels: failed to get channels from database")
	}
	return channels, nil
}

func (a *WhatsappApp) CreateWhatsappChannel(channelId string) (*model.WhatsappChannel, error) {
	channel, err := a.store.CreateWhatsappChannel(channelId)
	if err != nil {
		return nil, errors.Wrap(err, "CreateWhatsappChannel: failed to create channel")
	}
	return channel, nil
}
