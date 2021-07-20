import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

export default class LinkTooltip extends React.PureComponent {
    static propTypes = {
        href: PropTypes.string.isRequired,
        theme: PropTypes.object.isRequired,
    }

    render() {
        if (!this.props.href.includes('www.test.com')) {
            return null;
        }

        const style = getStyle(this.props.theme);
        return (
            <div
                style={style.configuration}
            >
                <i
                    style={iconStyles}
                    className='icon fa fa-plug'
                />
                <FormattedMessage
                    id='tooltip.message'
                    defaultMessage='This is a custom tooltip from the Demo Plugin'
                />
            </div>
        );
    }
}

const getStyle = (theme) => ({
    configuration: {
        borderRadius: '4px',
        boxShadow: 'rgba(61, 60, 64, 0.1) 0px 17px 50px 0px, rgba(61, 60, 64, 0.1) 0px 12px 15px 0px',
        fontSize: '14px',
        marginTop: '10px',
        padding: '10px 15px 15px',
        border: `1px solid ${hexToRGB(theme.centerChannelColor, '0.16')}`,
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
});

const iconStyles = {
    paddingRight: '5px',
};

export const hexToRGB = (hex, alpha) => {
    const r = parseInt(hex.slice(1, 3), 16);
    const g = parseInt(hex.slice(3, 5), 16);
    const b = parseInt(hex.slice(5, 7), 16);
    if (alpha) {
        return 'rgba(' + r + ', ' + g + ', ' + b + ', ' + alpha + ')';
    }
    return 'rgb(' + r + ', ' + g + ', ' + b + ')';
};
