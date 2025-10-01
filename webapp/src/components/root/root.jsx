import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import ChartsDialog from 'components/ChartsDialog';

const Root = ({visible, close, theme, subMenu}) => {
    if (!visible) {
        return null;
    }

    const style = getStyle(theme);

    const isNode = React.isValidElement(subMenu);

    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
            <div
                style={style.modal}
                onClick={(e) => e.stopPropagation()}
            >
                {isNode ? (
                    subMenu
                ) : (
                    <ChartsDialog/>
                )}
            </div>
        </div>
    );
};

Root.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
    subMenu: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
};

const getStyle = (theme) => ({
    backdrop: {
        position: 'fixed',
        display: 'flex',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.45)',
        zIndex: 2000,
        alignItems: 'center',
        justifyContent: 'center',
        padding: 16,
    },
    modal: {
        maxHeight: '90vh',
        width: 'min(92vw, 920px)',
        padding: 16,
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
        borderRadius: 8,
        boxShadow: '0 10px 25px rgba(0,0,0,0.2)',
        overflowY: 'auto',
    },
});

export default Root;
