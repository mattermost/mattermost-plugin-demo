// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import ContentCopyIcon from '@mattermost/compass-icons/components/content-copy';
import PencilOutlineIcon from '@mattermost/compass-icons/components/pencil-outline';
import TrashCanOutlineIcon from '@mattermost/compass-icons/components/trash-can-outline';
import React, {useLayoutEffect, useRef, useState} from 'react';
import {createPortal} from 'react-dom';

import {Button} from '@mattermost/compass-ui/components/button';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {MenuItem} from '@mattermost/compass-ui/components/menu-item';
import {Modal} from '@mattermost/compass-ui/components/modal';
import {PopoverMenu} from '@mattermost/compass-ui/components/popover-menu';
import {TextInput} from '@mattermost/compass-ui/components/text-input';
import {useOutsideClose} from '@mattermost/compass-ui/hooks/use-outside-close';

import ModalLayer from '../../../modal_layer/modal_layer';

export function LiveModal() {
    const [open, setOpen] = useState(false);
    const close = () => setOpen(false);

    return (
        <>
            <Button onClick={() => setOpen(true)}>{'Open live modal'}</Button>
            <ModalLayer
                open={open}
                onClose={close}
            >
                <Modal
                    open={open}
                    size='small'
                    title='Rename channel'
                    closeOnEscape={false}
                    onClose={close}
                    footer={(
                        <>
                            <Button
                                emphasis='tertiary'
                                onClick={close}
                            >
                                {'Cancel'}
                            </Button>
                            <Button onClick={close}>{'Save'}</Button>
                        </>
                    )}
                >
                    <TextInput
                        label='Channel name'
                        defaultValue='Town Square'
                    />
                </Modal>
            </ModalLayer>
        </>
    );
}

const TABBABLE = 'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

function focusBeside(element: HTMLElement, backwards: boolean) {
    const tabbable = Array.from(document.querySelectorAll<HTMLElement>(TABBABLE)).
        filter((el) => el === element || (!el.closest('[inert]') && el.getClientRects().length > 0));
    tabbable[tabbable.indexOf(element) + (backwards ? -1 : 1)]?.focus();
}

const MENU_ITEMS = [
    {label: 'Edit', glyph: <PencilOutlineIcon/>},
    {label: 'Copy text', glyph: <ContentCopyIcon/>},
    {label: 'Delete', glyph: <TrashCanOutlineIcon/>},
];

export function LiveMenu() {
    const [open, setOpen] = useState(false);
    const [position, setPosition] = useState({top: 0, left: 0});
    const triggerRef = useRef<HTMLButtonElement>(null);
    const menuRef = useRef<HTMLDivElement>(null);

    const close = (restoreFocus = true) => {
        setOpen(false);
        if (restoreFocus) {
            triggerRef.current?.focus();
        }
    };
    useOutsideClose(triggerRef, open, () => close(false), menuRef);

    useLayoutEffect(() => {
        const trigger = triggerRef.current;
        if (!open || !trigger) {
            return undefined;
        }
        const place = () => {
            const rect = trigger.getBoundingClientRect();
            setPosition({top: rect.bottom + 4, left: rect.left});
        };
        place();
        menuRef.current?.querySelector<HTMLElement>('[role="menuitem"]')?.focus();

        window.addEventListener('resize', place);
        window.addEventListener('scroll', place, true);
        return () => {
            window.removeEventListener('resize', place);
            window.removeEventListener('scroll', place, true);
        };
    }, [open]);

    const onKeyDown = (e: React.KeyboardEvent) => {
        const items = Array.from(menuRef.current?.querySelectorAll<HTMLElement>('[role="menuitem"]') ?? []);
        const index = items.indexOf(document.activeElement as HTMLElement);
        if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
            e.preventDefault();
            const step = e.key === 'ArrowDown' ? 1 : -1;
            items[(index + step + items.length) % items.length]?.focus();
        } else if (e.key === 'Escape') {
            e.preventDefault();
            close();
        } else if (e.key === 'Tab' && triggerRef.current) {
            e.preventDefault();
            close(false);
            focusBeside(triggerRef.current, e.shiftKey);
        }
    };

    return (
        <>
            <Button
                ref={triggerRef}
                emphasis='secondary'
                aria-haspopup='menu'
                aria-expanded={open}
                onClick={() => (open ? close() : setOpen(true))}
            >
                {'Open live menu'}
            </Button>
            {open && createPortal(
                <div
                    ref={menuRef}
                    className='DemoCompassShowcase__menuLayer'
                    style={position}
                    onKeyDown={onKeyDown}
                >
                    <PopoverMenu
                        role='menu'
                        aria-label='Message actions'
                    >
                        {MENU_ITEMS.map(({label, glyph}) => (
                            <MenuItem
                                key={label}
                                role='menuitem'
                                tabIndex={-1}
                                label={label}
                                leadingVisual={<Icon glyph={glyph}/>}
                                destructive={label === 'Delete'}
                                onClick={() => close()}
                            />
                        ))}
                    </PopoverMenu>
                </div>,
                document.body,
            )}
        </>
    );
}
