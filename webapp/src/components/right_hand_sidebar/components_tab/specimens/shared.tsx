// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

type Props = {
    children: React.ReactNode;
    label?: string;
};

export function Row({children}: {children: React.ReactNode}) {
    return (
        <div className='DemoCompassShowcase__row'>
            {children}
        </div>
    );
}

export function Variant({children, label}: Props) {
    return (
        <div className='DemoCompassShowcase__variant'>
            {label && (
                <div className='DemoCompassShowcase__variantLabel'>
                    {label}
                </div>
            )}
            <div className='DemoCompassShowcase__variantLive'>
                {children}
            </div>
        </div>
    );
}

export function Note({children}: {children: React.ReactNode}) {
    return (
        <p className='DemoCompassShowcase__note'>
            {children}
        </p>
    );
}

export function Excerpt({children}: {children: React.ReactNode}) {
    return (
        <div className='DemoCompassShowcase__excerpt'>
            {children}
        </div>
    );
}

export function WidePreview({
    children,
    center,
    padded,
}: {
    children: React.ReactNode;
    center?: boolean;
    padded?: boolean;
}) {
    let className = 'DemoCompassShowcase__widePreview';
    if (padded) {
        className += ' DemoCompassShowcase__widePreview--padded';
    }
    if (center) {
        className += ' DemoCompassShowcase__widePreview--center';
    }

    return (
        <div className={className}>
            {children}
        </div>
    );
}
