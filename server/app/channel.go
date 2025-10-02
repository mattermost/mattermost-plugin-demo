package app

import (
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (a *WhatsappApp) CreateChannel(channelID string) (*model.Channel, error) {
	channel := &model.Channel{
		ChannelID: channelID,
	}
	channel.SetDefaults()
	if err := channel.IsValid(); err != nil {
		return nil, errors.Wrap(err, "CreateChannel: channel is not valid")
	}
	if err := a.store.CreateChannel(channel, channelID); err != nil {
		return nil, errors.Wrap(err, "CreateChannel: failed to create channel")
	}
	return channel, nil
}

func (a *WhatsappApp) GetChannelByID(channelID string) (*model.Channel, error) {
	channel, err := a.store.GetChannelByID(channelID)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannelByID: failed to get channel")
	}
	return channel, nil
}

func (a *WhatsappApp) UpdateChannel(channel *model.Channel) error {
	if err := channel.IsValid(); err != nil {
		return errors.Wrap(err, "UpdateChannel: channel is not valid")
	}
	if err := a.store.UpdateChannel(channel); err != nil {
		return errors.Wrap(err, "UpdateChannel: failed to update channel")
	}
	return nil
}

func (a *WhatsappApp) GetChannels() ([]*model.Channel, error) {
	channels, err := a.store.GetChannels()
	if err != nil {
		return nil, errors.Wrap(err, "GetChannels: failed to get channels")
	}
	return channels, nil
}

func (a *WhatsappApp) DeleteChannel(channelID string) error {
	if err := a.store.DeleteChannel(channelID); err != nil {
		return errors.Wrap(err, "DeleteChannel: failed to delete channel")
	}
	return nil
}
