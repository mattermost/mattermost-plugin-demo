// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {type ReactNode, useEffect, useRef} from 'react';
import {createPortal} from 'react-dom';

import {modalOverlayStyles} from '@mattermost/compass-ui/components/modal';

import './modal_layer.css';

type Props = {
    onClose: () => void;
    children: ReactNode;
};

export default function ModalLayer({onClose, children}: Props) {
    const closeRef = useRef(onClose);
    closeRef.current = onClose;

    useEffect(() => {
        const {overflow} = document.body.style;
        document.body.style.overflow = 'hidden';
        const onEscape = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && !(e.target as Element).closest?.('[aria-expanded="true"], [role="listbox"]')) {
                e.stopPropagation();
                closeRef.current();
            }
        };
        window.addEventListener('keydown', onEscape, true);
        return () => {
            document.body.style.overflow = overflow;
            window.removeEventListener('keydown', onEscape, true);
        };
    }, []);

    return createPortal(
        <div
            className={`${modalOverlayStyles['modal-overlay']} DemoModalLayer`}
            onMouseDown={(e) => e.target === e.currentTarget && onClose()}
        >
            {children}
        </div>,
        document.body,
    );
}
