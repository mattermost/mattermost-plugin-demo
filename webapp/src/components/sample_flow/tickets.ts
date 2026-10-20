// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {useSyncExternalStore} from 'react';

import messages from './messages';

export type Ticket = {
    id: string;
    key: string;
    title: string;
    description: string;
    project: string;
    reporter: string;
    reported: string;
    labels: string[];
    status: string;
    priority: string;
    notifyWatchers: boolean;
};

export const STATUS_OPTIONS = [
    {value: 'todo', label: messages.statusTodo},
    {value: 'in_progress', label: messages.statusInProgress},
    {value: 'in_review', label: messages.statusInReview},
    {value: 'done', label: messages.statusDone},
];

export const PRIORITY_OPTIONS = [
    {value: 'high', label: messages.priorityHigh},
    {value: 'medium', label: messages.priorityMedium},
    {value: 'low', label: messages.priorityLow},
];

export const LABEL_OPTIONS = ['UX Needed', 'Groom', 'Mobile', 'Accessibility', 'Performance', 'Plugins'].
    map((label) => ({value: label, label}));

let tickets: Ticket[] = [
    {
        id: 'mm-38408',
        key: 'MM-38408',
        title: 'Mobile Webview: re-evaluate layout of the channel header',
        description: 'The header wraps on narrow screens and hides the call button.',
        project: 'Mattermost',
        reporter: 'Leonard Riley',
        reported: '2026-10-11',
        labels: ['UX Needed', 'Groom'],
        status: 'todo',
        priority: 'medium',
        notifyWatchers: true,
    },
    {
        id: 'mm-38512',
        key: 'MM-38512',
        title: 'Screen readers skip the unread filter in the sidebar',
        description: 'VoiceOver jumps from the team name to the first channel.',
        project: 'Mattermost',
        reporter: 'Aiko Tanaka',
        reported: '2026-10-16',
        labels: ['Accessibility'],
        status: 'in_progress',
        priority: 'high',
        notifyWatchers: true,
    },
    {
        id: 'mm-38590',
        key: 'MM-38590',
        title: 'Plugin settings page loads every plugin manifest',
        description: 'Opening one plugin\'s settings fetches all manifests, which is slow with many plugins.',
        project: 'Mattermost',
        reporter: 'Sam Okafor',
        reported: '2026-10-18',
        labels: ['Plugins', 'Performance'],
        status: 'in_review',
        priority: 'low',
        notifyWatchers: false,
    },
];

const listeners = new Set<() => void>();

function subscribe(listener: () => void) {
    listeners.add(listener);
    return () => {
        listeners.delete(listener);
    };
}

export function useTickets() {
    return useSyncExternalStore(subscribe, () => tickets);
}

export function saveTicket(next: Ticket) {
    tickets = tickets.map((ticket) => (ticket.id === next.id ? next : ticket));
    listeners.forEach((listener) => listener());
}
