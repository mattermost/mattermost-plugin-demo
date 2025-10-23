import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';
import {Toggle} from '@mattermost/design-system';

const Root = ({visible, close, theme, subMenu}) => {
    const [toggled, setToggle] = useState(false);
    const toggleToggle = useCallback(() => setToggle((val) => !val), []);
    if (!visible) {
        return null;
    }

    let extraContent = '';
    let extraContentTitle = '';
    if (subMenu) {
        extraContentTitle = (
            <FormattedMessage
                id='demo.triggeredby'
                defaultMessage='Element clicked in the menu: '
            />
        );
        extraContent = subMenu;
    }

    const style = getStyle(theme);

    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
            <h2>{'Toggles'}</h2>
            <label>
                {'Toggle'}
                <Toggle
                    onToggle={toggleToggle}
                    toggled={toggled}
                    onText='On!'
                    offText='Off!'
                />
            </label>
            <label>
                {'Disabled toggle'}
                <Toggle
                    onToggle={toggleToggle}
                    toggled={toggled}
                    disabled={true}
                    onText='On!'
                    offText='Off!'
                />
            </label>
            <label>
                {'Big toggle'}
                <Toggle
                    onToggle={toggleToggle}
                    toggled={toggled}
                    size='btn-lg'
                    onText='On!'
                    offText='Off!'
                />
            </label>

            <h2>{'Buttons'}</h2>

            <span>
                <Button
                    size='xs'
                >
                    {'Extra small'}
                </Button>
                <Button
                    size='sm'
                >
                    {'Small'}
                </Button>
            </span>
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
