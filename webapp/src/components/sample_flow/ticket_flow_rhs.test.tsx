// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {screen, within} from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {renderWithProviders} from 'tests/render';

import TicketFlowRhs from './ticket_flow_rhs';

describe('TicketFlowRhs', () => {
    it('edits a ticket, shares the change and returns focus to its row', async () => {
        renderWithProviders(<TicketFlowRhs pagePath='/team/plugin/tickets'/>);

        const row = screen.getByRole('button', {name: /Screen readers skip the unread filter/});
        await userEvent.click(row);

        const summary = screen.getByRole('textbox', {name: 'Summary'});
        expect(screen.getByRole('heading', {name: /Screen readers skip the unread filter/})).toBeInTheDocument();
        expect(screen.getByText(/MM-38512 • Mattermost/)).toBeInTheDocument();

        await userEvent.clear(summary);
        await userEvent.type(summary, 'Screen readers skip the filter');
        await userEvent.click(screen.getByRole('button', {name: 'Save'}));

        const updatedRow = screen.getByRole('button', {name: /Screen readers skip the filter/});
        expect(updatedRow).toHaveFocus();
        expect(within(updatedRow).getByText(/Reported Oct 16, 2026 by Aiko Tanaka/)).toBeInTheDocument();
    });

    it('discards edits on back and puts focus on the row it came from', async () => {
        renderWithProviders(<TicketFlowRhs pagePath='/team/plugin/tickets'/>);

        await userEvent.click(screen.getByRole('button', {name: /Plugin settings page loads every plugin manifest/}));
        await userEvent.type(screen.getByRole('textbox', {name: 'Summary'}), ' (draft)');
        await userEvent.click(screen.getByRole('button', {name: 'Back to tickets'}));

        expect(screen.queryByText(/\(draft\)/)).not.toBeInTheDocument();
        expect(screen.getByRole('button', {name: /Plugin settings page loads every plugin manifest/})).toHaveFocus();
    });
});
