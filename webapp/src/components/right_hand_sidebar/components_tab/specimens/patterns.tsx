// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {Button} from '@mattermost/compass-ui/components/button';
import {MenuItem} from '@mattermost/compass-ui/components/menu-item';
import {Modal} from '@mattermost/compass-ui/components/modal';
import {PopoverMenu} from '@mattermost/compass-ui/components/popover-menu';
import {ProfilePopover} from '@mattermost/compass-ui/components/profile-popover';
import {TourPoint} from '@mattermost/compass-ui/components/tour-point';

import {LiveMenu, LiveModal} from './live_overlays';
import {Note, WidePreview} from './shared';

export function ModalPreview() {
    return (
        <WidePreview>
            <Modal
                onClose={() => undefined}
                title='Rename channel'
                footer={(
                    <Button>{'Save'}</Button>
                )}
            >
                {'This surface is mount-controlled chrome. The host owns portal and focus trap.'}
            </Modal>
        </WidePreview>
    );
}

export function ModalDetail() {
    return (
        <>
            <Note>
                {'Modal handles focus, Tab and Escape. The plugin adds the portal, backdrop and scroll lock — try it:'}
            </Note>
            <LiveModal/>
            <Modal
                size='small'
                subtitle='Visible to everyone in this channel'
                title='Channel header'
                onBack={() => undefined}
                onClose={() => undefined}
                showBackButton={true}
                footer={(
                    <>
                        <Button emphasis='tertiary'>{'Cancel'}</Button>
                        <Button>{'Save'}</Button>
                    </>
                )}
            >
                {'Use Compass Modal as the visual shell; keep open/close in the product.'}
            </Modal>
        </>
    );
}

export function PopoverMenuPreview() {
    return (
        <PopoverMenu>
            <MenuItem label='Copy link'/>
            <MenuItem label='Mute channel'/>
        </PopoverMenu>
    );
}

export function PopoverMenuDetail() {
    return (
        <>
            <Note>
                {'PopoverMenu is only the surface. The plugin adds positioning, outside click, menu roles and arrow keys — try it:'}
            </Note>
            <LiveMenu/>
            <PopoverMenuStatic/>
        </>
    );
}

function PopoverMenuStatic() {
    return (
        <PopoverMenu>
            <MenuItem label='Edit'/>
            <MenuItem label='Copy text'/>
            <MenuItem
                destructive={true}
                label='Delete'
            />
        </PopoverMenu>
    );
}

export function ProfilePopoverPreview() {
    return (
        <ProfilePopover
            avatarAlt='Aiko Tan'
            avatarSrc=''
            name='Aiko Tan'
            username='aiko'
        />
    );
}

export function ProfilePopoverDetail() {
    return (
        <>
            <Note>
                {'Profile Popover is chrome-only — no portal or hover trigger in this plugin.'}
            </Note>
            <ProfilePopover
                avatarAlt='Leonard Riley'
                avatarSrc=''
                email='leonard@mattermost.com'
                name='Leonard Riley'
                title='Staff Software Engineer'
                username='leonard'
            />
        </>
    );
}

export function TourPointPreview() {
    return (
        <div className='DemoCompassShowcase__tourPointPreview'>
            <TourPoint title='Browse components'>
                {'Search or filter by category, then open a card for live variants.'}
            </TourPoint>
        </div>
    );
}

export function TourPointDetail() {
    return (
        <div className='DemoCompassShowcase__tourPointDetail'>
            <TourPoint
                pointerPosition='top-left'
                primaryAction={{label: 'Next'}}
                progress={{pages: 3, activePage: 1}}
                title='Compass UI'
            >
                {'This plugin is a production-style consumer of the published npm package.'}
            </TourPoint>
        </div>
    );
}
