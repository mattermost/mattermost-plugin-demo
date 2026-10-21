// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import PowerPlugOutlineIcon from '@mattermost/compass-icons/components/power-plug-outline';
import {Icon} from '@mattermost/compass-ui/components/icon';

import './demo_chrome.css';

export const PlugIcon = ({size = '16'}) => (
    <Icon
        glyph={<PowerPlugOutlineIcon/>}
        size={size}
    />
);

PlugIcon.propTypes = {
    size: PropTypes.string,
};

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
