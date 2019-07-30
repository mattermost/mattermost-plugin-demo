import React from 'react';
import PropTypes from 'prop-types';

export default class ChannelSubHeader extends React.PureComponent {
    static propTypes = {
        channelId: PropTypes.object.isRequired,
        theme: PropTypes.object.isRequired,
    }

    render() {
        const {channelId, theme} = this.props;
        const style = getStyle(theme);

        if (channelId !== 'k6nkyqnubbgobfwy6iab83ruio') {
            return null;
        }

        return (
            <div style={style.container}>
                {'This is a channel sub header from the demo plugin.'}
                <iframe
                    src='https://grafana.test.cloud.mattermost.com/d-solo/MQQu0iHZz/mattermost-test-cloud?orgId=1&panelId=2'
                    width='450'
                    height='200'
                    frameBorder='0'
                />
            </div>
        );
    }
}

const getStyle = (theme) => ({
    container: {
        padding: '5px',
        height: '30px',
        textAlign: 'center',
        borderTop: '1px solid rgba(61, 60, 64, 0.2)',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
});
