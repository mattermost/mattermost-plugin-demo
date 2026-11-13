// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {MessageHeader} from '@mattermost/compass-ui/components/message-header';
import {MessageSeparator} from '@mattermost/compass-ui/components/message-separator';
import {ThreadFooter} from '@mattermost/compass-ui/components/thread-footer';
import {ThreadListItem} from '@mattermost/compass-ui/components/thread-list-item';

import {Note, WidePreview} from './shared';

export function MessageHeaderPreview() {
    return (
        <MessageHeader
            timestamp='10:04 AM'
            username='Aiko Tan'
        />
    );
}

export function MessageHeaderDetail() {
    return (
        <MessageHeader
            isBot={true}
            timestamp='9:41 AM'
            username='Welcome Bot'
        />
    );
}

export function MessageSeparatorPreview() {
    return (
        <MessageSeparator label='Today'/>
    );
}

export function MessageSeparatorDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <MessageSeparator
                label='Today'
                type='date'
            />
            <MessageSeparator type='new-messages'/>
            <MessageSeparator
                label='4 replies'
                type='reply-count'
            />
        </div>
    );
}

export function ThreadFooterPreview() {
    return (
        <ThreadFooter replyCount={6}/>
    );
}

export function ThreadFooterDetail() {
    return (
        <ThreadFooter
            badge='mention'
            following={true}
            lastReplyTime='2 min ago'
            mentionCount={2}
            replyCount={11}
        />
    );
}

export function ThreadListItemPreview() {
    return (
        <WidePreview>
            <ThreadListItem
                authorName='Arjun Patel'
                previewText='Let’s keep the specimens compact so they fit the RHS.'
            />
        </WidePreview>
    );
}

export function ThreadListItemDetail() {
    return (
        <>
            <Note>
                {'Thread list rows reveal overflow actions on hover and focus-within.'}
            </Note>
            <ThreadListItem
                authorName='Sofia Bauer'
                badge='unread'
                channelLabel='UX Design'
                previewText='Search and filter chrome should be Compass components.'
            />
        </>
    );
}
