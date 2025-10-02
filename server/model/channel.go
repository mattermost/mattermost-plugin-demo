package model

import (
	"errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/utils"
)

type Channel struct {
	ID            string `json:"id"`
	ChannelID     string `json:"channel_id"`
	PhoneNumber   string `json:"phone_number"`
	PhoneNumberID string `json:"phone_number_id"`
}

func (c *Channel) SetDefaults() {
	if c.ID == "" {
		c.ID = utils.NewID()
	}

}

func (c *Channel) IsValid() error {
	if c.ID == "" {
		return errors.New("channel ID cannot be empty")
	}
	if c.ChannelID == "" {
		return errors.New("channel ChannelID cannot be empty")
	}
	if c.PhoneNumber == "" {
		return errors.New("channel PhoneNumber cannot be empty")
	}
	if c.PhoneNumberID == "" {
		return errors.New("channel PhoneNumberID cannot be empty")
	}
	return nil
}
