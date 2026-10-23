// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {ChannelSidebarItem} from '@mattermost/compass-ui/components/channel-sidebar-item';
import {MenuItem} from '@mattermost/compass-ui/components/menu-item';

import {Note} from './shared';

export function MenuItemPreview() {
    return (
        <MenuItem label='Copy link'/>
    );
}

export function MenuItemDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <MenuItem label='Mark as unread'/>
            <MenuItem
                label='Mute channel'
                secondaryLabel='Until tomorrow'
            />
            <MenuItem
                destructive={true}
                label='Leave channel'
            />
        </div>
    );
}

export function ChannelSidebarItemPreview() {
    return (
        <ChannelSidebarItem name='off-topic'/>
    );
}

export function ChannelSidebarItemDetail() {
    return (
        <>
            <Note>
                {'Wide channel lists are excerpted here so they stay usable at ~400px.'}
            </Note>
            <div className='DemoCompassShowcase__stack'>
                <ChannelSidebarItem
                    active={true}
                    name='town-square'
                />
                <ChannelSidebarItem
                    name='ux-design'
                    status='unread'
                />
                <ChannelSidebarItem
                    leadingVisual='private'
                    mentionCount={4}
                    name='design-leads'
                    status='mention'
                />
            </div>
        </>
    );
}
