// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {defineMessages, useIntl} from 'react-intl';

import {AdminPanel} from '@mattermost/compass-ui/components/admin-panel';
import {Button} from '@mattermost/compass-ui/components/button';
import {TextArea} from '@mattermost/compass-ui/components/text-area';

import {id as pluginId} from '../../manifest';

import './admin_settings.css';

const messages = defineMessages({
    title: {id: 'admin.customSettings.title', defaultMessage: 'Custom Settings'},
    subtitle: {
        id: 'admin.customSettings.subtitle',
        defaultMessage: 'Plugin-owned System Console controls registered with registerAdminConsoleCustomSetting, built with Compass UI.',
    },
    messageLabel: {id: 'admin.customSettings.messageLabel', defaultMessage: 'Secret message'},
    show: {id: 'admin.customSettings.show', defaultMessage: 'Show Secret Message'},
    hide: {id: 'admin.customSettings.hide', defaultMessage: 'Hide Secret Message'},
    sliderLabel: {id: 'admin.customSettings.sliderLabel', defaultMessage: 'Demo custom setting'},
});

type Props = {
    id: string;
    value: string;
    disabled: boolean;
    setByEnv: boolean;
    config?: {
        PluginSettings?: {
            Plugins?: Record<string, Record<string, unknown>>;
        };
    };
    onChange: (id: string, value: unknown) => void;
    registerSaveAction: (action: () => Promise<{error?: {message?: string}}>) => void;
    unRegisterSaveAction: (action: () => Promise<{error?: {message?: string}}>) => void;
};

function pluginSetting(config: Props['config'], key: string) {
    return config?.PluginSettings?.Plugins?.[pluginId]?.[key];
}

export default function SecretMessageSetting({
    id,
    value,
    disabled,
    setByEnv,
    config,
    onChange,
    registerSaveAction,
    unRegisterSaveAction,
}: Props) {
    const {formatMessage} = useIntl();
    const [showMessage, setShowMessage] = useState(false);
    const [slider, setSlider] = useState(() => Number(pluginSetting(config, 'CustomSetting') ?? 0));

    useEffect(() => {
        const handleSave = async () => ({error: undefined as {message?: string} | undefined});
        registerSaveAction(handleSave);
        return () => unRegisterSaveAction(handleSave);
    }, [registerSaveAction, unRegisterSaveAction]);

    return (
        <AdminPanel
            className='DemoAdminSetting'
            title={formatMessage(messages.title)}
            subtitle={formatMessage(messages.subtitle)}
        >
            <div className='DemoAdminSetting__body'>
                {showMessage && (
                    <TextArea
                        label={formatMessage(messages.messageLabel)}
                        rows={5}
                        value={value ?? ''}
                        disabled={disabled || setByEnv}
                        onChange={(e) => onChange(id, e.target.value)}
                    />
                )}
                <Button
                    emphasis='tertiary'
                    disabled={disabled}
                    onClick={() => setShowMessage((open) => !open)}
                >
                    {formatMessage(showMessage ? messages.hide : messages.show)}
                </Button>
                <div className='DemoAdminSetting__sliderBlock'>
                    <p className='DemoAdminSetting__sliderValue'>
                        {`${formatMessage(messages.sliderLabel)}: ${slider}`}
                    </p>
                    <input
                        className='DemoAdminSetting__slider'
                        type='range'
                        min={0}
                        max={10}
                        disabled={disabled}
                        value={slider}
                        aria-label={formatMessage(messages.sliderLabel)}
                        style={{'--demo-slider-pct': `${(slider / 10) * 100}%`} as React.CSSProperties}
                        onChange={(e) => {
                            const next = Number(e.target.value);
                            setSlider(next);
                            onChange('CustomSetting', next);
                        }}
                    />
                </div>
            </div>
        </AdminPanel>
    );
}
