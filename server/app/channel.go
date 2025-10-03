package app

import (
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (app *WhatsappApp) CreateChannel(channelID string, phoneNumber string, phoneNumberId string) (*model.Channel, error) {
	channel := &model.Channel{
		ChannelID:     channelID,
		PhoneNumber:   phoneNumber,
		PhoneNumberID: phoneNumberId,
	}
	channel.SetDefaults()
	if err := channel.IsValid(); err != nil {
		return nil, errors.Wrap(err, "CreateChannel: channel is not valid")
	}
	if err := app.store.CreateChannel(channel); err != nil {
		return nil, errors.Wrap(err, "CreateChannel: failed to create channel")
	}
	return channel, nil
}

func (app *WhatsappApp) GetChannelByID(channelID string) (*model.Channel, error) {
	channel, err := app.store.GetChannelByID(channelID)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannelByID: failed to get channel")
	}
	return channel, nil
}

func (app *WhatsappApp) UpdateChannel(channel *model.Channel) error {
	if err := channel.IsValid(); err != nil {
		return errors.Wrap(err, "UpdateChannel: channel is not valid")
	}
	if err := app.store.UpdateChannel(channel); err != nil {
		return errors.Wrap(err, "UpdateChannel: failed to update channel")
	}
	return nil
}

func (app *WhatsappApp) GetChannels() ([]*model.Channel, error) {
	channels, err := app.store.GetChannels()
	if err != nil {
		return nil, errors.Wrap(err, "GetChannels: failed to get channels")
	}
	return channels, nil
}

func (app *WhatsappApp) DeleteChannel(channelID string) error {
	if err := app.store.DeleteChannel(channelID); err != nil {
		return errors.Wrap(err, "DeleteChannel: failed to delete channel")
	}
	return nil
}
