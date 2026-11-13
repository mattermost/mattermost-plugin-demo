// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {type ReactNode} from 'react';
import {defineMessages, useIntl} from 'react-intl';

import {AdminPanel} from '@mattermost/compass-ui/components/admin-panel';

import './admin_settings.css';

const messages = defineMessages({
    title: {id: 'admin.customSettings.title', defaultMessage: 'Custom Settings'},
    subtitle: {
        id: 'admin.customSettings.subtitle',
        defaultMessage: 'Plugin-owned System Console controls registered with registerAdminConsoleCustomSection.',
    },
});

type Props = {
    settingsList: ReactNode[];
};

// Groups SecretMessage + CustomSetting in one AdminPanel via registerAdminConsoleCustomSection.
export default function CustomSettingsSection({settingsList}: Props) {
    const {formatMessage} = useIntl();

    return (
        <AdminPanel
            className='DemoAdminSetting'
            title={formatMessage(messages.title)}
            subtitle={formatMessage(messages.subtitle)}
        >
            <div className='DemoAdminSetting__body'>
                {settingsList}
            </div>
        </AdminPanel>
    );
}
