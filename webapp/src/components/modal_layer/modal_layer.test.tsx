// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {fireEvent, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {modalOverlayStyles} from '@mattermost/compass-ui/components/modal';

import {renderWithProviders} from 'tests/render';

import ModalLayer from './modal_layer';

const content = (
    <div role='dialog'>
        <button type='button'>{'Inside'}</button>
        <div
            role='combobox'
            aria-expanded='true'
            tabIndex={0}
        >
            {'Open select'}
        </div>
    </div>
);

function renderLayer(onClose = jest.fn()) {
    const result = renderWithProviders(
        <ModalLayer
            open={true}
            onClose={onClose}
        >
            {content}
        </ModalLayer>,
    );
    const close = () => result.rerender(
        <ModalLayer
            open={false}
            onClose={onClose}
        >
            {content}
        </ModalLayer>,
    );
    return {...result, onClose, close};
}

const overlay = () => screen.getByRole('dialog').parentElement as HTMLElement;
const expectUnmounted = () => expect(screen.queryByRole('dialog')).not.toBeInTheDocument();

describe('ModalLayer', () => {
    it('renders nothing while closed', () => {
        renderWithProviders(
            <ModalLayer
                open={false}
                onClose={jest.fn()}
            >
                {content}
            </ModalLayer>,
        );
        expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('portals to the body and locks page scroll until it unmounts', () => {
        document.body.style.overflow = 'auto';
        const {unmount, container} = renderLayer();

        expect(container).toBeEmptyDOMElement();
        expect(overlay().parentElement).toBe(document.body);
        expect(document.body.style.overflow).toBe('hidden');

        unmount();
        expect(document.body.style.overflow).toBe('auto');
    });

    it('stays mounted with the exiting class until the exit finishes', async () => {
        document.body.style.overflow = 'auto';
        const {close, onClose} = renderLayer();

        close();
        expect(overlay()).toHaveClass(modalOverlayStyles['modal-overlay--exiting']);
        expect(document.body.style.overflow).toBe('hidden');

        fireEvent.mouseDown(overlay());
        await userEvent.keyboard('{Escape}');
        expect(onClose).not.toHaveBeenCalled();

        await waitFor(expectUnmounted);
        expect(document.body.style.overflow).toBe('auto');
    });

    it('closes on Escape and stops the host from handling the same keypress', async () => {
        const hostHandler = jest.fn();
        document.addEventListener('keydown', hostHandler);
        const {onClose} = renderLayer();

        screen.getByRole('button', {name: 'Inside'}).focus();
        await userEvent.keyboard('{Escape}');

        expect(onClose).toHaveBeenCalledTimes(1);
        expect(hostHandler).not.toHaveBeenCalled();
        document.removeEventListener('keydown', hostHandler);
    });

    it('leaves Escape to an open select inside the modal', async () => {
        const {onClose} = renderLayer();

        screen.getByRole('combobox').focus();
        await userEvent.keyboard('{Escape}');

        expect(onClose).not.toHaveBeenCalled();
    });

    it('closes on a backdrop press but not on a press inside the modal', () => {
        const {onClose} = renderLayer();

        fireEvent.mouseDown(screen.getByRole('button', {name: 'Inside'}));
        expect(onClose).not.toHaveBeenCalled();

        fireEvent.mouseDown(overlay());
        expect(onClose).toHaveBeenCalledTimes(1);
    });
});
