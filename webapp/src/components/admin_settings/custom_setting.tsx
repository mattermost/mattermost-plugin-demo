// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {defineMessages, useIntl} from 'react-intl';

import './admin_settings.css';

const messages = defineMessages({
    sliderLabel: {id: 'admin.customSettings.sliderLabel', defaultMessage: 'Demo custom setting'},
});

type Props = {
    id: string;
    value?: number | string | null;
    disabled: boolean;
    onChange: (id: string, value: unknown) => void;
    setSaveNeeded: () => void;
};

export default function CustomSetting({id, value, disabled, onChange, setSaveNeeded}: Props) {
    const {formatMessage} = useIntl();
    const [slider, setSlider] = useState(() => Number(value ?? 0));

    return (
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
                    onChange(id, next);
                    setSaveNeeded();
                }}
            />
        </div>
    );
}
