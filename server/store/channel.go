package store

import (
	"database/sql"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	"github.com/pkg/errors"
)

func (s *SQLStore) channelColumns() []string {
	return []string{
		"id",
		"channel_id",
	}
}

func (s *SQLStore) GetChannels() ([]*model.Channel, error) {
	rows, err := s.getQueryBuilder().
		Select(s.channelColumns()...).
		From(s.tablePrefix + "channel").
		Query()

	if err != nil {
		return nil, err
	}
	channels, err := s.ChannelsFromRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannels: failed to map channel rows to channels")
	}
	return channels, nil
}

func (s *SQLStore) ChannelsFromRows(rows *sql.Rows) ([]*model.Channel, error) {
	channels := []*model.Channel{}
	for rows.Next() {
		var channel model.Channel
		err := rows.Scan(
			&channel.ID,
			&channel.ChannelID,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &channel)
	}
	return channels, nil
}
