import React from 'react';
import PropTypes from 'prop-types';

export default class UserAttributes extends React.PureComponent {
    static propTypes = {
        hide: PropTypes.func.isRequired,
    }

    onClick = () => {
        this.props.hide();
    }

    render() {
        return (
            <div>
                <a onClick={this.onClick}>
                    {'Demo Plugin: User Attributes'}
                </a>
            </div>
        );
    }
}

