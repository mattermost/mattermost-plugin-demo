import React from 'react';
import PropTypes from 'prop-types';

const SystemWideSettingModal = ({visible, close, theme, setting}) => {
    if (!visible) {
        return null;
    }

    const style = getStyle(theme);

    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
            <div style={style.modal}>
                <h3>{ 'System wide setting has changed.' }</h3>
                <br/>
                <br/>
                <p>{ 'New setting is: ' } { setting }</p>
            </div>
        </div>
    );
};

SystemWideSettingModal.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
    setting: PropTypes.string.isRequired,
};

const getStyle = (theme) => ({
    backdrop: {
        position: 'absolute',
        display: 'flex',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.50)',
        zIndex: 2000,
        alignItems: 'center',
        justifyContent: 'center',
    },
    modal: {
        height: '250px',
        width: '400px',
        padding: '1em',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
});

export default SystemWideSettingModal;
