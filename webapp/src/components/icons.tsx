// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import PowerPlugOutlineIcon from '@mattermost/compass-icons/components/power-plug-outline';
import {Icon, type IconSize} from '@mattermost/compass-ui/components/icon';

import './demo_chrome.css';

type PlugIconProps = {
    size?: IconSize;
};

export const PlugIcon = ({size = '16'}: PlugIconProps) => (
    <Icon
        glyph={<PowerPlugOutlineIcon/>}
        size={size}
    />
);

export const MainMenuMobileIcon = () => (
    <PlugIcon size='16'/>
);

// Use a div, not a span — host app-bar styles force a white background on spans.
export const ChannelHeaderButtonIcon = () => (
    <div className='DemoChannelHeaderButtonIcon'>
        <PlugIcon size='16'/>
    </div>
);

export const FileUploadMethodIcon = () => (
    <PlugIcon size='16'/>
);
