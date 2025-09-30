import React from 'react';
import PropTypes from 'prop-types';

const WHATSAPP_GREEN = '#25D366';

export default class LeftSidebarHeader extends React.PureComponent {
    static propTypes = {
        enabled: PropTypes.bool.isRequired,
    }

    render() {
        const {enabled} = this.props;

        const iconStyle = {
            display: 'inline-block',
            marginRight: '8px',
            marginLeft: '5px',
            color: enabled ? WHATSAPP_GREEN : 'rgba(255,255,255,0.4)',
        };

        const statusTextStyle = {
            fontSize: '12px',
            fontWeight: 'bold',
            marginLeft: '5px',
            color: enabled ? WHATSAPP_GREEN : 'rgba(255,255,255,0.4)',
            textShadow: enabled ? `0 0 3px ${WHATSAPP_GREEN}, 0 0 1px #000` : 'none',
        };

        const headerStyle = {
            margin: '0',
            padding: '7px 16px 7px 19px',
            color: 'rgba(255,255,255,0.7)',
            backgroundColor: 'transparent',
        };

        return (
            <div style={headerStyle}>
                <i
                    className='icon fa fa-brands fa-whatsapp'
                    style={iconStyle}
                />
                {'WhatsApp:'}
                <span style={statusTextStyle}>
                    {enabled ? 'ON' : 'OFF'}
                </span>
            </div>
        );
    }
}
