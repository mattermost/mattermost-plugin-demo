import React from 'react';
import {FormattedMessage} from 'react-intl';
import {useSelector} from 'react-redux';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

export default function ChannelInfoSection() {
    const channel = useSelector(getCurrentChannel);

    if (!channel) {
        return null;
    }

    return (
        <div>
            <div>
                <strong>
                    <FormattedMessage
                        id='channel_settings.info_section.name'
                        defaultMessage='Channel:'
                    />
                </strong>
                {' '}
                {channel.display_name}
            </div>
            <div>
                <strong>
                    <FormattedMessage
                        id='channel_settings.info_section.type'
                        defaultMessage='Type:'
                    />
                </strong>
                {' '}
                {channel.type}
            </div>
            <p style={{marginTop: '8px', opacity: 0.75}}>
                <FormattedMessage
                    id='channel_settings.info_section.note'
                    defaultMessage='This section is rendered entirely by the demo plugin to showcase a custom settings section.'
                />
            </p>
        </div>
    );
}
