// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';

import {PaginationDots} from '@mattermost/compass-ui/components/pagination-dots';
import {ProgressBar} from '@mattermost/compass-ui/components/progress-bar';
import {Spinner} from '@mattermost/compass-ui/components/spinner';

import {Row, Variant} from './shared';

export function SpinnerPreview() {
    return (
        <Spinner size='24'/>
    );
}

export function SpinnerDetail() {
    return (
        <Row>
            <Spinner size='16'/>
            <Spinner size='24'/>
            <Spinner size='32'/>
        </Row>
    );
}

export function ProgressBarPreview() {
    return (
        <div style={{width: 112}}>
            <ProgressBar value={64}/>
        </div>
    );
}

export function ProgressBarDetail() {
    return (
        <>
            <Variant label='Default'>
                <ProgressBar value={40}/>
            </Variant>
            <Variant label='Semantic'>
                <ProgressBar
                    semanticColors={true}
                    size='small'
                    value={92}
                />
            </Variant>
        </>
    );
}

export function PaginationDotsPreview() {
    return (
        <PaginationDots
            activePage={2}
            pages={4}
        />
    );
}

export function PaginationDotsDetail() {
    const [activePage, setActivePage] = useState(1);

    return (
        <PaginationDots
            activePage={activePage}
            pages={5}
            onPageChange={setActivePage}
        />
    );
}
