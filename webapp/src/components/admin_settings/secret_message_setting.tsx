// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {defineMessages, useIntl} from 'react-intl';

import {Button} from '@mattermost/compass-ui/components/button';
import {TextArea} from '@mattermost/compass-ui/components/text-area';

import './admin_settings.css';

const messages = defineMessages({
    messageLabel: {id: 'admin.customSettings.messageLabel', defaultMessage: 'Secret message'},
    show: {id: 'admin.customSettings.show', defaultMessage: 'Show Secret Message'},
    hide: {id: 'admin.customSettings.hide', defaultMessage: 'Hide Secret Message'},
});

type Props = {
    id: string;
    value: string;
    disabled: boolean;
    setByEnv: boolean;
    onChange: (id: string, value: unknown) => void;
    registerSaveAction: (action: () => Promise<{error?: {message?: string}}>) => void;
    unRegisterSaveAction: (action: () => Promise<{error?: {message?: string}}>) => void;
};

export default function SecretMessageSetting({
    id,
    value,
    disabled,
    setByEnv,
    onChange,
    registerSaveAction,
    unRegisterSaveAction,
}: Props) {
    const {formatMessage} = useIntl();
    const [showMessage, setShowMessage] = useState(false);

    useEffect(() => {
        const handleSave = async () => ({error: undefined as {message?: string} | undefined});
        registerSaveAction(handleSave);
        return () => unRegisterSaveAction(handleSave);
    }, [registerSaveAction, unRegisterSaveAction]);

    return (
        <>
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
        </>
    );
}
