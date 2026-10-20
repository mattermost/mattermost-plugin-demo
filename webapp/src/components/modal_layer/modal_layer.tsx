// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {type ReactNode, useEffect, useRef} from 'react';
import {createPortal} from 'react-dom';

import {modalOverlayStyles} from '@mattermost/compass-ui/components/modal';
import {usePopoverTransition} from '@mattermost/compass-ui/hooks/use-popover-transition';

import './modal_layer.css';

type Props = {
    open: boolean;
    onClose: () => void;
    children: ReactNode;
};

export default function ModalLayer({open, onClose, children}: Props) {
    const {mounted} = usePopoverTransition(open);
    const closeRef = useRef(onClose);
    closeRef.current = onClose;

    useEffect(() => {
        if (!mounted) {
            return undefined;
        }
        const {overflow} = document.body.style;
        document.body.style.overflow = 'hidden';
        return () => {
            document.body.style.overflow = overflow;
        };
    }, [mounted]);

    useEffect(() => {
        if (!open) {
            return undefined;
        }
        const onEscape = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && !(e.target as Element).closest?.('[aria-expanded="true"], [role="listbox"]')) {
                e.stopPropagation();
                closeRef.current();
            }
        };
        window.addEventListener('keydown', onEscape, true);
        return () => window.removeEventListener('keydown', onEscape, true);
    }, [open]);

    if (!mounted) {
        return null;
    }

    const className = [modalOverlayStyles['modal-overlay'], !open && modalOverlayStyles['modal-overlay--exiting'], 'DemoModalLayer'].filter(Boolean).join(' ');

    return createPortal(
        <div
            className={className}
            onMouseDown={(e) => open && e.target === e.currentTarget && onClose()}
        >
            {children}
        </div>,
        document.body,
    );
}
