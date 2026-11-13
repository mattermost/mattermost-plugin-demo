// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import PaletteOutlineIcon from '@mattermost/compass-icons/components/palette-outline';
import {Emoji} from '@mattermost/compass-ui/components/emoji';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {Illustration} from '@mattermost/compass-ui/components/illustration';
import {TeamAvatar} from '@mattermost/compass-ui/components/team-avatar';
import {UserAvatar} from '@mattermost/compass-ui/components/user-avatar';
import {UserAvatarGroup} from '@mattermost/compass-ui/components/user-avatar-group';

import {Row, Variant} from './shared';

export function IconPreview() {
    return (
        <Icon
            glyph={<PaletteOutlineIcon/>}
            size='24'
        />
    );
}

export function IconDetail() {
    return (
        <Row>
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='16'
            />
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='24'
            />
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='32'
            />
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='52'
            />
        </Row>
    );
}

export function EmojiPreview() {
    return (
        <Emoji emoji='👍'/>
    );
}

export function EmojiDetail() {
    return (
        <Row>
            <Emoji
                emoji='🎉'
                size='16'
            />
            <Emoji
                emoji='🚀'
                size='24'
            />
            <Emoji
                emoji='✅'
                size='32'
            />
        </Row>
    );
}

export function UserAvatarPreview() {
    return (
        <UserAvatar
            alt='Leonard Riley'
            name='Leonard Riley'
            size='40'
            status={true}
        />
    );
}

export function UserAvatarDetail() {
    return (
        <Variant label='Fallback sizes'>
            <Row>
                <UserAvatar
                    alt='Aiko Tan'
                    name='Aiko Tan'
                    size='24'
                />
                <UserAvatar
                    alt='Arjun Patel'
                    name='Arjun Patel'
                    size='40'
                    status={true}
                />
                <UserAvatar
                    alt='Sofia Bauer'
                    name='Sofia Bauer'
                    size='64'
                />
            </Row>
        </Variant>
    );
}

export function UserAvatarGroupPreview() {
    return (
        <UserAvatarGroup
            avatars={[
                {key: 'a', name: 'Aiko Tan'},
                {key: 'b', name: 'Arjun Patel'},
                {key: 'c', name: 'Sofia Bauer'},
                {key: 'd', name: 'Marco Rinaldi'},
            ]}
        />
    );
}

export function UserAvatarGroupDetail() {
    return (
        <UserAvatarGroup
            avatars={[
                {key: 'a', name: 'Aiko Tan'},
                {key: 'b', name: 'Arjun Patel'},
                {key: 'c', name: 'Sofia Bauer'},
                {key: 'd', name: 'Marco Rinaldi'},
                {key: 'e', name: 'Leonard Riley'},
            ]}
            max={3}
            size='32'
        />
    );
}

export function TeamAvatarPreview() {
    return (
        <TeamAvatar
            alt='Contributors'
            initials='Co'
            size='40'
        />
    );
}

export function TeamAvatarDetail() {
    return (
        <div className='DemoCompassShowcase__sidebarWell'>
            <TeamAvatar
                alt='Contributors'
                initials='Co'
                size='32'
            />
            <TeamAvatar
                alt='Design'
                badge={3}
                initials='De'
                size='40'
                state='active'
            />
        </div>
    );
}

export function IllustrationMark() {
    return (
        <svg
            aria-hidden='true'
            viewBox='0 0 160 96'
            xmlns='http://www.w3.org/2000/svg'
        >
            <ellipse
                cx='80'
                cy='80'
                fill='rgba(var(--center-channel-color-rgb), 0.08)'
                rx='40'
                ry='7'
            />
            <circle
                cx='80'
                cy='44'
                fill='none'
                r='26'
                stroke='rgba(var(--center-channel-color-rgb), 0.28)'
                strokeWidth='2'
            />
            <circle
                cx='80'
                cy='44'
                fill='rgba(var(--button-bg-rgb), 0.16)'
                r='14'
                stroke='var(--button-bg)'
                strokeWidth='2'
            />
            <rect
                fill='rgba(var(--center-channel-color-rgb), 0.45)'
                height='7'
                rx='1'
                width='2'
                x='79'
                y='14'
            />
            <rect
                fill='rgba(var(--center-channel-color-rgb), 0.45)'
                height='7'
                rx='1'
                width='2'
                x='79'
                y='67'
            />
            <rect
                fill='rgba(var(--center-channel-color-rgb), 0.45)'
                height='2'
                rx='1'
                width='7'
                x='50'
                y='43'
            />
            <rect
                fill='rgba(var(--center-channel-color-rgb), 0.45)'
                height='2'
                rx='1'
                width='7'
                x='103'
                y='43'
            />
            <polygon
                fill='var(--button-bg)'
                points='80,24 87,48 80,44 73,48'
            />
            <polygon
                fill='rgba(var(--center-channel-color-rgb), 0.32)'
                points='80,64 87,40 80,44 73,40'
            />
            <circle
                cx='80'
                cy='44'
                fill='var(--button-bg)'
                r='3.5'
            />
        </svg>
    );
}

export function IllustrationPreview() {
    return (
        <Illustration
            aria-label='Sample illustration'
            height='48px'
            width='80px'
        >
            <IllustrationMark/>
        </Illustration>
    );
}

export function IllustrationDetail() {
    return (
        <Illustration
            aria-label='Sample illustration'
            height='96px'
            width='160px'
        >
            <IllustrationMark/>
        </Illustration>
    );
}
