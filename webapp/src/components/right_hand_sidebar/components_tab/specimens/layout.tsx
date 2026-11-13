// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';

import {Divider} from '@mattermost/compass-ui/components/divider';
import {EmptyState} from '@mattermost/compass-ui/components/empty-state';
import {Scrollbar} from '@mattermost/compass-ui/components/scrollbar';
import {Tabs} from '@mattermost/compass-ui/components/tabs';

import {IllustrationMark} from './images';
import {Note, Variant, WidePreview} from './shared';

const emptyStateIllustration = {
    'aria-label': 'Sample illustration',
    children: <IllustrationMark/>,
    height: '96px',
    width: '160px',
};

export function DividerPreview() {
    return (
        <div
            className='DemoCompassShowcase__stack'
            style={{width: 88}}
        >
            <span>{'Above'}</span>
            <Divider/>
            <span>{'Below'}</span>
        </div>
    );
}

export function DividerDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <span>{'Above'}</span>
            <Divider/>
            <span>{'Below'}</span>
        </div>
    );
}

export function EmptyStatePreview() {
    return (
        <WidePreview>
            <EmptyState
                description='Saved messages will show up here.'
                illustration={emptyStateIllustration}
                title='No saved messages'
            />
        </WidePreview>
    );
}

export function EmptyStateDetail() {
    return (
        <EmptyState
            action={{children: 'Clear search'}}
            description='Try a different name or category.'
            illustration={emptyStateIllustration}
            title='No components match'
        />
    );
}

export function ScrollbarPreview() {
    return (
        <Scrollbar
            alwaysVisible={true}
            className='DemoCompassShowcase__scrollbarPreview'
            style={{height: 80, width: 64}}
        >
            <div
                className='DemoCompassShowcase__scrollbarDemo'
                style={{height: 160}}
            >
                <div className='DemoCompassShowcase__scrollbarDemoLine'/>
                <div className='DemoCompassShowcase__scrollbarDemoLine'/>
                <div className='DemoCompassShowcase__scrollbarDemoLine'/>
                <div className='DemoCompassShowcase__scrollbarDemoLine'/>
                <div className='DemoCompassShowcase__scrollbarDemoLine'/>
                <div className='DemoCompassShowcase__scrollbarDemoLine'/>
            </div>
        </Scrollbar>
    );
}

export function ScrollbarDetail() {
    return (
        <Scrollbar
            alwaysVisible={true}
            style={{height: 96}}
        >
            <div
                className='DemoCompassShowcase__scrollbarDemo'
                style={{height: 180}}
            >
                {'alwaysVisible keeps the thumb on screen while content overflows.'}
            </div>
        </Scrollbar>
    );
}

export function TabsPreview() {
    return (
        <Tabs
            activeKey='all'
            onChange={() => undefined}
            tabs={[
                {key: 'all', label: 'All', countBadge: 12},
                {key: 'actions', label: 'Actions', countBadge: 3},
            ]}
        />
    );
}

export function TabsDetail() {
    const [activeKey, setActiveKey] = useState('mentions');
    return (
        <>
            <Note>
                {'Tabs is the Compass pattern for filters on a list — the gallery uses it with countBadge.'}
            </Note>
            <Variant label='With counts and unread'>
                <Tabs
                    activeKey={activeKey}
                    onChange={setActiveKey}
                    tabs={[
                        {key: 'mentions', label: 'Mentions', countBadge: 4},
                        {key: 'threads', label: 'Threads', unreadBadge: true},
                        {key: 'saved', label: 'Saved'},
                    ]}
                />
            </Variant>
        </>
    );
}
