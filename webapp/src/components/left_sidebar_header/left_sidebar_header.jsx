import React from 'react';
import PropTypes from 'prop-types';

const WHATSAPP_GREEN = '#25D366';
const WHATSAPP_DARK_GREEN = '#075E54';

export default class LeftSidebarHeader extends React.PureComponent {
    static propTypes = {
        enabled: PropTypes.bool.isRequired,
    }

    render() {
        const {enabled} = this.props;

        const iconStyle = {
            display: 'inline-block',
            margin: '0 7px 0 1px',
            color: enabled ? WHATSAPP_GREEN : '#3f4350b8',
        };

        const headerStyle = {
            margin: '.5em 0 .5em',
            padding: '4px 12px 4px 15px',
            borderRadius: '4px',
            fontWeight: '600',
            transition: 'all 0.3s ease',
        };

        let backgroundStyle;
        if (enabled) {
            backgroundStyle = {
                backgroundColor: WHATSAPP_DARK_GREEN,
                color: '#FFFFFF',
            };
        } else {
            backgroundStyle = {
                backgroundColor: 'rgba(177,187,208,0.72)',
                color: '#3f4350',
            };
        }

        const finalStyle = {...headerStyle, ...backgroundStyle};

        return (
            <div style={finalStyle}>
                <i
                    className='icon fa fa-brands fa-whatsapp'
                    style={iconStyle}
                />
                {'WhatsApp:'}
                {' '}
                {enabled ?
                    <span style={{color: WHATSAPP_GREEN, marginLeft: '5px'}}>
                        <b>{'ON'}</b>
                    </span> :
                    <span style={{color: '#3f4350', marginLeft: '5px'}}>
                        <b>{'OFF'}</b>
                    </span>
                }
            </div>
        );
    }
}
