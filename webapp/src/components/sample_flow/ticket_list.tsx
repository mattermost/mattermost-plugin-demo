// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useRef} from 'react';
import {FormattedDate, FormattedMessage} from 'react-intl';

import {Tag} from '@mattermost/compass-ui/components/tag';

import jiraStory from './jira_story.svg';
import messages from './messages';
import type {Ticket} from './tickets';

export function Reported({ticket}: {ticket: Ticket}) {
    return (
        <FormattedMessage
            {...messages.reported}
            values={{
                reporter: ticket.reporter,
                date: (
                    <FormattedDate
                        key='date'
                        value={ticket.reported}
                        dateStyle='medium'
                        timeZone='UTC'
                    />
                ),
            }}
        />
    );
}

type Props = {
    tickets: Ticket[];
    onSelect: (id: string) => void;
    selectedId?: string;
    focusId?: string;
};

export default function TicketList({tickets, onSelect, selectedId, focusId}: Props) {
    const listRef = useRef<HTMLUListElement>(null);

    useEffect(() => {
        if (focusId) {
            listRef.current?.querySelector<HTMLButtonElement>(`[data-ticket-id="${focusId}"]`)?.focus();
        }
    }, [focusId]);

    return (
        <ul
            ref={listRef}
            className='DemoTicketList'
        >
            {tickets.map((ticket) => (
                <li key={ticket.id}>
                    <button
                        type='button'
                        className='DemoTicketList__item'
                        data-ticket-id={ticket.id}
                        aria-current={ticket.id === selectedId || undefined}
                        onClick={() => onSelect(ticket.id)}
                    >
                        <span className='DemoTicketList__title'>
                            <img
                                src={jiraStory}
                                alt=''
                                width={14}
                                height={14}
                            />
                            <span className='DemoTicketList__summary'>{ticket.title}</span>
                        </span>
                        {ticket.labels.length > 0 && (
                            <span className='DemoTicketList__tags'>
                                {ticket.labels.map((label) => (
                                    <Tag
                                        key={label}
                                        label={label}
                                        size='x-small'
                                        casing='title-case'
                                    />
                                ))}
                            </span>
                        )}
                        <span className='DemoTicketList__meta'>
                            <span>{`${ticket.key} • ${ticket.project}`}</span>
                            <span className='DemoTicketList__muted'><Reported ticket={ticket}/></span>
                        </span>
                    </button>
                </li>
            ))}
        </ul>
    );
}
