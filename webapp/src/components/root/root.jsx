import React, {useCallback, useEffect, useState} from 'react';
import PropTypes from 'prop-types';

import {
    Button,
    CheckInput,
    IconButton,
    Input,
    PasswordInput,
} from '@mattermost/design-system';

const Root = ({visible, close, theme, subMenu}) => {
    const [checked1, toggleCheck1] = useToggle();
    const [checked2, toggleCheck2] = useToggle();

    const [toggled1, toggleToggle1] = useToggle();
    const [toggled2, toggleToggle2] = useToggle();
    const [toggled3, toggleToggle3] = useToggle();
    const [toggled4, toggleToggle4] = useToggle();

    useEffect(() => {
        const closeOnEscape = (e) => {
            console.log(e.key);
            if (e.key === 'Escape') {
                close();
            }
        };

        if (visible) {
            document.addEventListener('keydown', closeOnEscape);

            return () => {
                document.removeEventListener('keydown', closeOnEscape);
            }
        } else {
            return () => {};
        }
    }, [visible, close]);

    if (!visible) {
        return null;
    }

    const style = getStyle(theme);

    return (
        <div style={style.backdrop}>
            <div style={style.modal}>
                <div style={style.modalContent}>

                    <h2>{'Buttons'}</h2>

                    <div style={style.row}>
                        <Button size='xs'>
                            {'Extra small'}
                        </Button>
                        <Button size='sm'>
                            {'Small'}
                        </Button>
                        <Button size='md'>
                            {'Medium'}
                        </Button>
                        <Button
                            size='lg'
                            destructive={true}
                        >
                            {'Large'}
                        </Button>
                    </div>

                    <h2>{'Checkbox'}</h2>
                        <div style={style.row}>
                        <label>
                            {'Checkbox'}
                            <CheckInput
                                checked={checked1}
                                onChange={toggleCheck1}
                            />
                        </label>
                        <label>
                            {'Disabled checkbox'}
                            <CheckInput
                                checked={checked2}
                                onChange={toggleCheck2}
                                disabled={true}
                            />
                        </label>
                    </div>

                    <h2>{'Input'}</h2>
                    <Input
                        type='text'
                        placeholder='Enter text...'
                    />
                    <PasswordInput
                        placeholder='Enter password...'
                    />

                    <h2>{'Icon button'}</h2>
                    <div style={style.row}>
                        <IconButton
                            icon={<span style={{fontSize: 20}}>{'≤'}</span>}
                            size='lg'
                            inverted={false}
                            onClick={toggleToggle1}
                            toggled={toggled1}
                        />
                        <div style={{backgroundColor: 'darkblue'}}>
                            <IconButton
                                icon={<span style={{fontSize: 20}}>{'≤'}</span>}
                                size='lg'
                                inverted={true}
                                onClick={toggleToggle2}
                                toggled={toggled2}
                            />
                        </div>
                        <IconButton
                            icon={<span style={{fontSize: 20}}>{'☹'}</span>}
                            size='lg'
                            destructive={true}
                            onClick={toggleToggle3}
                            toggled={toggled3}
                        />
                        <div style={{backgroundColor: 'darkblue'}}>
                            <IconButton
                                icon={<span style={{fontSize: 20}}>{'☹'}</span>}
                                size='lg'
                                destructive={true}
                                inverted={true}
                                onClick={toggleToggle4}
                                toggled={toggled4}
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

function useToggle() {
    const [toggled, setToggle] = useState(false);
    const toggleToggle = useCallback(() => setToggle((val) => !val), []);

    return [toggled, toggleToggle];
}

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
        height: '80vh',
        width: '80vw',
        overflowY: 'scroll',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
        display: 'flex',
    },
    modalContent: {
        padding: '1em',
        width: '60%',
    },
    row: {
        display: 'flex',
        flexDirection: 'row',
        'gap': '16px',
    }
});

export default Root;
