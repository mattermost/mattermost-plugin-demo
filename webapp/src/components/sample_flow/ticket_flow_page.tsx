// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useId, useRef} from 'react';
import {FormattedMessage, defineMessages, useIntl} from 'react-intl';
import {useHistory, useLocation, useRouteMatch} from 'react-router-dom';

import {Button} from '@mattermost/compass-ui/components/button';
import {EmptyState} from '@mattermost/compass-ui/components/empty-state';

import PluginPage from '../plugin_page/plugin_page';

import sharedMessages from './messages';
import TicketForm from './ticket_form';
import TicketList, {Reported} from './ticket_list';
import {useTickets} from './tickets';

import './sample_flow.css';

const messages = defineMessages({
    description: {id: 'ticketFlow.page.description', defaultMessage: 'The sample flow at full width. The same list and form as the sidebar and the modal, laid out side by side.'},
    details: {id: 'ticketFlow.page.details', defaultMessage: 'Ticket details'},
    emptyTitle: {id: 'ticketFlow.page.emptyTitle', defaultMessage: 'No ticket selected'},
    emptyDescription: {id: 'ticketFlow.page.emptyDescription', defaultMessage: 'Pick a ticket from the list to edit it.'},
});

export default function TicketFlowPage() {
    const {formatMessage} = useIntl();
    const tickets = useTickets();
    const formId = useId();
    const history = useHistory();
    const {search} = useLocation();
    const {url} = useRouteMatch();
    const selectedId = new URLSearchParams(search).get('ticket') ?? undefined;
    const selected = tickets.find((ticket) => ticket.id === selectedId);

    const listRef = useRef<HTMLElement>(null);

    const open = (id: string) => history.push(`${url}?ticket=${id}`);
    const close = () => {
        history.push(url);
        listRef.current?.querySelector<HTMLButtonElement>(`[data-ticket-id="${selectedId}"]`)?.focus();
    };

    return (
        <PluginPage
            title={formatMessage(sharedMessages.tickets)}
            description={formatMessage(messages.description)}
        >
            <div className='DemoTicketFlowPage'>
                <nav
                    ref={listRef}
                    aria-label={formatMessage(sharedMessages.tickets)}
                    className='DemoTicketFlowPage__list'
                >
                    <TicketList
                        tickets={tickets}
                        onSelect={open}
                        selectedId={selected?.id}
                    />
                </nav>
                <section
                    className='DemoTicketFlowPage__detail'
                    aria-labelledby={selected ? `${formId}-title` : undefined}
                    aria-label={selected ? undefined : formatMessage(messages.details)}
                >
                    {selected ? (
                        <>
                            <header>
                                <h2
                                    id={`${formId}-title`}
                                    className='DemoTicketFlowPage__title'
                                >
                                    {selected.key}
                                </h2>
                                <p className='DemoTicketFlowPage__meta'>
                                    {`${selected.project} • `}
                                    <Reported ticket={selected}/>
                                </p>
                            </header>
                            <TicketForm
                                key={selected.id}
                                ticket={selected}
                                formId={formId}
                                onSaved={close}
                            />
                            <div className='DemoTicketFlow__actions'>
                                <Button
                                    type='submit'
                                    form={formId}
                                >
                                    <FormattedMessage {...sharedMessages.save}/>
                                </Button>
                                <Button
                                    emphasis='tertiary'
                                    onClick={close}
                                >
                                    <FormattedMessage {...sharedMessages.cancel}/>
                                </Button>
                            </div>
                        </>
                    ) : (
                        <EmptyState
                            title={formatMessage(messages.emptyTitle)}
                            description={formatMessage(messages.emptyDescription)}
                        />
                    )}
                </section>
            </div>
        </PluginPage>
    );
}
