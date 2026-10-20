// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {fireEvent, render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ModalLayer from './modal_layer';

function renderLayer(onClose = jest.fn()) {
    const result = render(
        <ModalLayer onClose={onClose}>
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
        </ModalLayer>,
    );
    return {...result, onClose};
}

describe('ModalLayer', () => {
    it('portals to the body and locks page scroll until it unmounts', () => {
        document.body.style.overflow = 'auto';
        const {unmount, container} = renderLayer();

        expect(container).toBeEmptyDOMElement();
        expect(screen.getByRole('dialog').closest('.DemoModalLayer')?.parentElement).toBe(document.body);
        expect(document.body.style.overflow).toBe('hidden');

        unmount();
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

        fireEvent.mouseDown(document.querySelector('.DemoModalLayer') as Element);
        expect(onClose).toHaveBeenCalledTimes(1);
    });
});
