// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {Button} from '@mattermost/compass-ui/components/button';
import {FileInfo} from '@mattermost/types/files';

import '../demo_chrome.css';

type Props = {
    fileInfo: FileInfo;
    onModalDismissed: () => void;
};

export default function FilePreviewOverride({fileInfo, onModalDismissed}: Props) {
    return (
        <div className='DemoFilePreviewOverride'>
            <h3 className='DemoFilePreviewOverride__title'>{fileInfo.name}</h3>
            <Button
                emphasis='primary'
                onClick={onModalDismissed}
            >
                {'Close'}
            </Button>
        </div>
    );
}
