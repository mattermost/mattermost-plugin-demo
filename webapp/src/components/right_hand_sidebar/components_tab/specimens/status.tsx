// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {MentionBadge} from '@mattermost/compass-ui/components/mention-badge';
import {ShortcutTag, ShortcutTagGroup} from '@mattermost/compass-ui/components/shortcut-tag';
import {StatusBadge} from '@mattermost/compass-ui/components/status-badge';
import {Tag} from '@mattermost/compass-ui/components/tag';
import {UnreadBadge} from '@mattermost/compass-ui/components/unread-badge';

import {Row} from './shared';

export function TagPreview() {
    return (
        <Tag label='Guest'/>
    );
}

export function TagDetail() {
    return (
        <Row>
            <Tag
                label='Default'
                type='default'
            />
            <Tag
                label='Info'
                type='info'
            />
            <Tag
                label='Success'
                type='success'
            />
            <Tag
                label='Warning'
                type='warning'
            />
            <Tag
                label='Danger'
                type='danger'
            />
        </Row>
    );
}

export function ShortcutTagPreview() {
    return (
        <ShortcutTagGroup labels={['⌘', 'K']}/>
    );
}

export function ShortcutTagDetail() {
    return (
        <Row>
            <ShortcutTag label='⌘'/>
            <ShortcutTag label='Shift'/>
            <ShortcutTag label='F'/>
        </Row>
    );
}

export function MentionBadgePreview() {
    return (
        <MentionBadge count={3}/>
    );
}

export function MentionBadgeDetail() {
    return (
        <Row>
            <MentionBadge
                count={1}
                location='sidebar'
            />
            <MentionBadge
                count={12}
                location='channel'
                size='medium'
            />
            <MentionBadge count={120}/>
        </Row>
    );
}

export function StatusBadgePreview() {
    return (
        <StatusBadge
            size='medium'
            status='online'
        />
    );
}

export function StatusBadgeDetail() {
    return (
        <Row>
            <StatusBadge status='online'/>
            <StatusBadge status='away'/>
            <StatusBadge status='do-not-disturb'/>
            <StatusBadge status='offline'/>
        </Row>
    );
}

export function UnreadBadgePreview() {
    return (
        <UnreadBadge/>
    );
}

export function UnreadBadgeDetail() {
    return (
        <Row>
            <UnreadBadge size='8'/>
            <UnreadBadge
                context='icon-button'
                size='6'
            />
        </Row>
    );
}
