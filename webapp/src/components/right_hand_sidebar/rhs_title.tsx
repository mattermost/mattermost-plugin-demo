// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage, defineMessages, useIntl} from 'react-intl';
import {useSelector} from 'react-redux';

import {Tag} from '@mattermost/compass-ui/components/tag';

import {isEnabled} from '../../selectors';

import './rhs_view.css';

const messages = defineMessages({
    hooksEnabled: {id: 'rhs.hooks.enabled', defaultMessage: 'Hooks enabled'},
    hooksDisabled: {id: 'rhs.hooks.disabled', defaultMessage: 'Hooks disabled'},
});

// Passed to registerRightHandSidebarComponent as title so the host header can
// show the hooks status beside the plugin name.
export default function RHSTitle() {
    const {formatMessage} = useIntl();
    const enabled = useSelector(isEnabled);

    return (
        <>
            <FormattedMessage
                id='plugin.name'
                defaultMessage='Demo Plugin'
            />
            <Tag
                className='demo-rhs__header-tag'
                label={formatMessage(enabled ? messages.hooksEnabled : messages.hooksDisabled)}
                size='x-small'
                type={enabled ? 'success' : 'warning'}
            />
        </>
    );
}
