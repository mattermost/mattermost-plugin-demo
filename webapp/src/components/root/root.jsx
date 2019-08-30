import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import {menu} from '../../sub_menu';

const Root = ({visible, close, theme, subMenuId}) => {
    if (!visible) {
        return null;
    }

    let extraContent = '';
    let extraContentTitle = '';
    if (subMenuId) {
        extraContentTitle = (
            <FormattedMessage
                id='demo.triggeredby'
                defaultMessage='Element clicked in the menu: '
            />
        );
        const items = [];
        items.push(menu);
        items.push(menu.subMenu);
        menu.subMenu.forEach((i) => items.push(i.subMenu));
        extraContent = items.flat().find((i) => i.id === subMenuId).text;

    }

    const style = getStyle(theme);

    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
            <div style={style.modal}>
                <FormattedMessage
                    id='root.triggered'
                    defaultMessage='You have triggered the root component of the demo plugin.'
                />
                <br/>
                <br/>
                <FormattedMessage
                    id='root.clicktoclose'
                    defaultMessage='Click anywhere to close.'
                />
                <br/>
                <br/>
                <FormattedMessage
                    id='demo.testintl'
                    defaultMessage='This is the default string'
                />
                <br/>
                <br/>
                {extraContentTitle}
                {extraContent}
            </div>
        </div>
    );
};

Root.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
    subMenuId: PropTypes.string,
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

export default Root;
