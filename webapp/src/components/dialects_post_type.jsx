import React from 'react';
import PropTypes from 'prop-types';

const {formatText, messageHtmlToComponent} = window.PostUtils;

export default class PostType extends React.PureComponent {
    static propTypes = {
        post: PropTypes.object.isRequired,
        theme: PropTypes.object.isRequired,
    }

    render() {
        const post = {...this.props.post};
        const message = post.message + '...' || '';

        const formattedText = messageHtmlToComponent(formatText(message));

        return (
            <div>
                {formattedText}
            </div>
        );
    }
}
