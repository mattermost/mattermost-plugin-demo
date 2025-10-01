package store

import (
	"fmt"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	_ "github.com/lib/pq"
	mm_model "github.com/mattermost/mattermost/server/public/model"
)

func (S SQLStore) GetWhatsappChannels() ([]model.WhatsappChannel, error) {
	query := "SELECT id, channel_id, phone_number, phone_number_id FROM whatsapp_plugin_channel"
	var channels []model.WhatsappChannel
	rows, err := S.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var channel model.WhatsappChannel
		var id string
		err := rows.Scan(
			&id,
			&channel.ChannelId,
			&channel.PhoneNumber,
			&channel.PhoneNumberId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (S SQLStore) CreateWhatsappChannel(channelId string) (*model.WhatsappChannel, error) {
	newId := mm_model.NewId()
	channel := &model.WhatsappChannel{
		ChannelId:     channelId,
		PhoneNumber:   "",
		PhoneNumberId: "",
	}
	query := "INSERT INTO whatsapp_plugin_channel (id, channel_id, phone_number, phone_number_id) VALUES ($1, $2, $3, $4)"
	_, err := S.db.Exec(query, newId, channel.ChannelId, channel.PhoneNumber, channel.PhoneNumberId)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	return channel, nil
}
