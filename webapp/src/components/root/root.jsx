import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import {menu} from '../../sub_menu';

const Root = ({visible, close, theme, subMenuId}) => {
    if (!visible) {
        return null;
    }

    let extraContent = '';
    if (subMenuId) {
        if (subMenuId === menu.id) {
            extraContent = menu.text;
        } else {
            let menuSearch = menu.subMenu.find((s) => s.id === subMenuId);
            if (menuSearch) {
                extraContent = menuSearch.text;
            } else {
                menu.subMenu.forEach((sm) => {
                    menuSearch = sm.subMenu.find((s) => s.id === subMenuId);
                    if (menuSearch) {
                        extraContent = menuSearch.text;
                    }
                });
            }
        }
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
