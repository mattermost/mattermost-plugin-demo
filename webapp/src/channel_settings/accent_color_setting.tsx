import React, {useEffect, useRef, useState} from 'react';
import {FormattedMessage} from 'react-intl';
import {useSelector} from 'react-redux';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import {fetchSchemaValues} from './client';

const SWATCHES = ['#1e88e5', '#43a047', '#e53935', '#8e24aa', '#fb8c00'];

const ACCENT_COLOR_SETTING = 'accentColor';

type Props = {
    informChange: (name: string, value: string) => void;
};

// The host treats custom settings as uncontrolled: a "Reset" clears the host's
// tracked value but will NOT visually revert this swatch. Acceptable for a demo.
export default function AccentColorSetting({informChange}: Props) {
    const channelId = useSelector(getCurrentChannelId);
    const [value, setValue] = useState<string>(SWATCHES[0]);

    // The host may pass a new informChange each render; a ref keeps the hydrate
    // effect depending only on the channel id so it doesn't refetch/loop.
    const informChangeRef = useRef(informChange);
    informChangeRef.current = informChange;

    useEffect(() => {
        let active = true;

        if (channelId) {
            fetchSchemaValues(channelId).then((values) => {
                const loaded = values[ACCENT_COLOR_SETTING];
                if (active && loaded) {
                    setValue(loaded);

                    // Sync the host's tracked value with what we hydrated.
                    informChangeRef.current(ACCENT_COLOR_SETTING, loaded);
                }
            }).catch(() => {
                // Plugin owns its own error surfacing; fall back to the default swatch.
            });
        }

        return () => {
            active = false;
        };
    }, [channelId]);

    const onPick = (color: string) => {
        setValue(color);
        informChange(ACCENT_COLOR_SETTING, color);
    };

    return (
        <div>
            <label>
                <FormattedMessage
                    id='channel_settings.accent_color.title'
                    defaultMessage='Accent color'
                />
            </label>
            <div style={{display: 'flex', gap: '8px', marginTop: '8px'}}>
                {SWATCHES.map((color) => (
                    <button
                        key={color}
                        type='button'
                        onClick={() => onPick(color)}
                        aria-pressed={value === color}
                        style={{
                            width: '28px',
                            height: '28px',
                            borderRadius: '50%',
                            background: color,
                            border: value === color ? '3px solid var(--center-channel-color)' : '1px solid rgba(var(--center-channel-color-rgb), 0.2)',
                            cursor: 'pointer',
                        }}
                    />
                ))}
            </div>
        </div>
    );
}
