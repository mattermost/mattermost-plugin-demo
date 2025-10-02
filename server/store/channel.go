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

func (s *SQLStore) CreateChannel(channel *model.Channel, channelId string) error {
	channel.SetDefaults()
	channel.ChannelID = channelId
	if err := channel.IsValid(); err != nil {
		return err
	}

	_, err := s.getQueryBuilder().
		Insert(s.tablePrefix+"channel").
		Columns(s.channelColumns()...).
		Values(
			channel.ID,
			channel.ChannelID,
		).
		Exec()
	if err != nil {
		return errors.Wrap(err, "SQLStore.CreateChannel failed to insert channel into database")
	}
	return nil

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

func (s *SQLStore) GetChannelByID(id string) (*model.Channel, error) {
	row := s.getQueryBuilder().
		Select(s.channelColumns()...).
		From(s.tablePrefix+"channel").
		Where("id = ?", id).
		QueryRow()

	var channel model.Channel
	if err := row.Scan(&channel.ID, &channel.ChannelID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("channel not found")
		}
		return nil, errors.Wrap(err, "GetChannelByID: failed to scan channel")
	}
	return &channel, nil
}

func (s *SQLStore) UpdateChannel(channel *model.Channel) error {
	if err := channel.IsValid(); err != nil {
		return errors.Wrap(err, "UpdateChannel: invalid channel")
	}

	_, err := s.getQueryBuilder().
		Update(s.tablePrefix+"channel").
		Set("channel_id", channel.ChannelID).
		Where("id = ?", channel.ID).
		Exec()

	if err != nil {
		return errors.Wrap(err, "UpdateChannel: failed to update channel in database")
	}
	return nil
}

func (s *SQLStore) DeleteChannel(id string) error {
	result, err := s.getQueryBuilder().
		Delete(s.tablePrefix+"channel").
		Where("id = ?", id).
		Exec()

	if err != nil {
		return errors.Wrap(err, "DeleteChannel: failed to delete channel from database")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "DeleteChannel: failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("DeleteChannel: channel not found")
	}
	return nil
}
