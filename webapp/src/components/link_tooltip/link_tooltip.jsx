import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

export default class LinkTooltip extends React.PureComponent {
    static propTypes = {
        href: PropTypes.string.isRequired,
        theme: PropTypes.object.isRequired,
    }

    render() {
        if (!this.props.href.includes('example.com')) {
            return null;
        }

        const style = getStyle(this.props.theme);
        return (
            <div style={style.container}>
                <div style={style.header}>
                    <i
                        style={iconStyles}
                        className='icon fa fa-plug'
                    />
                    <FormattedMessage
                        id='tooltip.message'
                        defaultMessage='This is a custom tooltip from the Demo Plugin'
                    />
                </div>
                <div style={style.body}>
                    <a
                        href={this.props.href}
                        target='_blank'
                        rel='noopener noreferrer'
                        style={style.titleLink}
                        data-testid='demo-tooltip-title-link'
                    >
                        <FormattedMessage
                            id='tooltip.title'
                            defaultMessage='Demo Link Preview'
                        />
                    </a>

                    <p style={style.sharedVia}>
                        <FormattedMessage
                            id='tooltip.sharedVia'
                            defaultMessage='Shared via {link}'
                            values={{
                                link: (
                                    <a
                                        href={this.props.href}
                                        target='_blank'
                                        rel='noopener noreferrer'
                                        data-testid='demo-tooltip-shared-via-link'
                                    >
                                        {'example.com'}
                                    </a>
                                ),
                            }}
                        />
                    </p>

                    <p style={style.description}>
                        <FormattedMessage
                            id='tooltip.description'
                            defaultMessage='This is a sample description. Visit {link} for more details.'
                            values={{
                                link: (
                                    <a
                                        href={this.props.href}
                                        target='_blank'
                                        rel='noopener noreferrer'
                                        data-testid='demo-tooltip-description-link'
                                    >
                                        {'the original page'}
                                    </a>
                                ),
                            }}
                        />
                    </p>
                </div>
            </div>
        );
    }
}

const getStyle = (theme) => ({
    container: {
        borderRadius: '4px',
        boxShadow: 'rgba(61, 60, 64, 0.1) 0px 17px 50px 0px, rgba(61, 60, 64, 0.1) 0px 12px 15px 0px',
        fontSize: '14px',
        marginTop: '10px',
        padding: '10px 15px 15px',
        border: `1px solid ${hexToRGB(theme.centerChannelColor, '0.16')}`,
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
    header: {
        marginBottom: '8px',
    },
    body: {
        paddingLeft: '4px',
    },
    titleLink: {
        display: 'block',
        fontWeight: '600',
        fontSize: '14px',
        textDecoration: 'none',
        color: theme.centerChannelColor,
    },
    sharedVia: {
        fontSize: '12px',
        marginTop: '4px',
        color: hexToRGB(theme.centerChannelColor, '0.64'),
    },
    description: {
        fontSize: '12px',
        marginTop: '8px',
        lineHeight: '1.25',
        color: theme.centerChannelColor,
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
