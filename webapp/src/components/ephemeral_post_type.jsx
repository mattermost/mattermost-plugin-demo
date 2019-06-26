import React from 'react';
import PropTypes from 'prop-types';

const {formatText, messageHtmlToComponent} = window.PostUtils;

export default class EphemeralPostType extends React.PureComponent {
    static propTypes = {
        post: PropTypes.object.isRequired,
        theme: PropTypes.object.isRequired,
    }

    render() {
        const style = getStyle(this.props.theme);
        const post = {...this.props.post};
        const message = post.message || '';

        const formattedText = messageHtmlToComponent(formatText(message));

        return (
            <div>
                <pre style={style.configuration}>
                    {formattedText}
                </pre>
            </div>
        );
    }
}

const getStyle = (theme) => ({
    configuration: {
        padding: '1em',
        color: theme.centerChannelBg,
        'font-style': 'oblique',
        backgroundColor: theme.centerChannelColor,
    },
});
