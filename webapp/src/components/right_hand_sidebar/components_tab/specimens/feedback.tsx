// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {ErrorMessage} from '@mattermost/compass-ui/components/error-message';
import {PopoverNotice} from '@mattermost/compass-ui/components/popover-notice';
import {SectionNotice} from '@mattermost/compass-ui/components/section-notice';
import {Toast} from '@mattermost/compass-ui/components/toast';
import {Tooltip} from '@mattermost/compass-ui/components/tooltip';

import {Note, Variant, WidePreview} from './shared';

export function ErrorMessagePreview() {
    return (
        <ErrorMessage message='Enter a valid channel URL.'/>
    );
}

export function ErrorMessageDetail() {
    return (
        <ErrorMessage message='This field is required.'/>
    );
}

export function SectionNoticePreview() {
    return (
        <WidePreview>
            <SectionNotice
                title='Invite teammates'
                description='Share this workspace with people who should see Compass UI.'
                type='info'
            />
        </WidePreview>
    );
}

export function SectionNoticeDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <SectionNotice
                description='Something needs attention.'
                title='Warning'
                type='warning'
            />
            <SectionNotice
                description='The change was saved.'
                title='Success'
                type='success'
                onDismiss={() => undefined}
            />
            <SectionNotice
                description='This action cannot be undone.'
                title='Danger'
                type='danger'
                primaryButtonLabel='Delete'
            />
        </div>
    );
}

export function ToastPreview() {
    return (
        <Toast
            message='Draft saved'
            type='success'
        />
    );
}

export function ToastDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <Toast
                message='Connecting to the server…'
                type='general'
            />
            <Toast
                actionLabel='Undo'
                message='Message sent'
                type='info'
            />
            <Toast
                message='Unable to save changes'
                type='danger'
                onDismiss={() => undefined}
            />
        </div>
    );
}

export function TooltipPreview() {
    return (
        <Tooltip label='Add a reaction'/>
    );
}

export function TooltipDetail() {
    return (
        <>
            <Note>
                {'Tooltip is chrome-only — the host owns hover, delay, and positioning.'}
            </Note>
            <Variant label='Label only'>
                <Tooltip label='Add a reaction'/>
            </Variant>
            <Variant label='Label and hint'>
                <Tooltip
                    hint='Search this channel'
                    label='Find'
                />
            </Variant>
            <Variant label='Label and shortcut'>
                <Tooltip
                    label='Find'
                    shortcutKeys={[{label: '⌘'}, {label: 'F'}]}
                />
            </Variant>
            <Variant label='Label, hint, and shortcut'>
                <Tooltip
                    hint='Search this channel'
                    label='Find'
                    shortcutKeys={[{label: '⌘'}, {label: 'F'}]}
                />
            </Variant>
            <Variant label='Arrow positions'>
                <div className='DemoCompassShowcase__stack'>
                    <Tooltip
                        arrow='top'
                        label='Top'
                    />
                    <Tooltip
                        arrow='right'
                        label='Right'
                    />
                    <Tooltip
                        arrow='bottom'
                        label='Bottom'
                    />
                    <Tooltip
                        arrow='left'
                        label='Left'
                    />
                </div>
            </Variant>
        </>
    );
}

export function PopoverNoticePreview() {
    return (
        <WidePreview>
            <PopoverNotice
                title='Try keyboard shortcuts'
                variant='info'
            >
                {'Press ⌘K to jump to a channel.'}
            </PopoverNotice>
        </WidePreview>
    );
}

export function PopoverNoticeDetail() {
    return (
        <PopoverNotice
            title='Enable notifications'
            variant='info'
            actions={[{label: 'Enable', emphasis: 'primary'}]}
        >
            {'Stay up to date when teammates mention you.'}
        </PopoverNotice>
    );
}
