// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {renderWithProviders} from 'tests/render';

import TicketForm from './ticket_form';
import * as store from './tickets';

const ticket: store.Ticket = {
    id: 'mm-1',
    key: 'MM-1',
    title: 'Original summary',
    description: '',
    project: 'Mattermost',
    reporter: 'Sam',
    reported: '2026-10-18',
    labels: [],
    status: 'todo',
    priority: 'low',
    notifyWatchers: false,
};

function renderForm() {
    const onSaved = jest.fn();
    const saveTicket = jest.spyOn(store, 'saveTicket').mockImplementation(() => undefined);
    renderWithProviders(
        <>
            <TicketForm
                ticket={ticket}
                formId='form'
                onSaved={onSaved}
            />
            <button
                type='submit'
                form='form'
            >
                {'Save'}
            </button>
        </>,
    );
    return {onSaved, saveTicket};
}

describe('TicketForm', () => {
    it('focuses the summary when it opens', () => {
        renderForm();
        expect(screen.getByRole('textbox', {name: 'Summary'})).toHaveFocus();
    });

    it('blocks an empty summary, announces the error and returns focus to the field', async () => {
        const {onSaved, saveTicket} = renderForm();
        const summary = screen.getByRole('textbox', {name: 'Summary'});

        await userEvent.clear(summary);
        await userEvent.click(screen.getByRole('button', {name: 'Save'}));

        expect(screen.getByText('There are errors in the form below')).toBeInTheDocument();
        expect(screen.getByRole('alert')).toHaveTextContent('Enter a summary');
        expect(summary).toHaveAttribute('aria-invalid', 'true');
        expect(summary).toHaveAccessibleDescription('Enter a summary');
        expect(summary).toHaveFocus();
        expect(saveTicket).not.toHaveBeenCalled();
        expect(onSaved).not.toHaveBeenCalled();
    });

    it('saves the trimmed draft from a submit button outside the form', async () => {
        const {onSaved, saveTicket} = renderForm();
        const summary = screen.getByRole('textbox', {name: 'Summary'});

        await userEvent.clear(summary);
        await userEvent.type(summary, '  New summary  ');
        await userEvent.click(screen.getByRole('radio', {name: 'High'}));
        await userEvent.click(screen.getByRole('button', {name: 'Save'}));

        expect(saveTicket).toHaveBeenCalledWith({...ticket, title: 'New summary', priority: 'high'});
        expect(onSaved).toHaveBeenCalledTimes(1);
    });
});
