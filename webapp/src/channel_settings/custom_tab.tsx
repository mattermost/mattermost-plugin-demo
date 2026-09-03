import React, {useCallback, useEffect, useMemo, useRef, useState} from 'react';
import {FormattedMessage, injectIntl, type InjectedIntlProps} from 'react-intl';

import {fetchCustomState, saveCustomState, type CustomState} from './client';
import type {ChannelSettingsTabBodyProps} from './types';

const EMPTY_STATE: CustomState = {note: '', pinGreeting: false};

// Inline styles can't target ::placeholder, so theme the textarea via a scoped
// class using the global Mattermost theme CSS variables (works on any theme).
const NOTE_TEXTAREA_CLASS = 'demo-channel-note-textarea';
const NOTE_TEXTAREA_STYLES = `
.${NOTE_TEXTAREA_CLASS} {
    background: var(--center-channel-bg);
    color: var(--center-channel-color);
    border: 1px solid rgba(var(--center-channel-color-rgb), 0.16);
    border-radius: 4px;
    padding: 6px 8px;
}
.${NOTE_TEXTAREA_CLASS}::placeholder {
    color: rgba(var(--center-channel-color-rgb), 0.64);
}
`;

function isDirty(a: CustomState, b: CustomState): boolean {
    return a.note !== b.note || a.pinGreeting !== b.pinGreeting;
}

// @types/react-intl is pinned to v2, which has no useIntl hook; injectIntl is
// the supported way to get a typed formatMessage for the placeholder attribute.
function ChannelSettingsCustomTab({channel, setUnsaved, registerHandlers, intl}: ChannelSettingsTabBodyProps & InjectedIntlProps) {
    const [baseline, setBaseline] = useState<CustomState>(EMPTY_STATE);
    const [state, setState] = useState<CustomState>(EMPTY_STATE);
    const [loading, setLoading] = useState(true);
    const [loadError, setLoadError] = useState(false);

    const dirty = useMemo(() => isDirty(state, baseline), [state, baseline]);

    useEffect(() => {
        setUnsaved(dirty);
    }, [dirty, setUnsaved]);

    useEffect(() => {
        let active = true;
        setLoading(true);
        setLoadError(false);
        fetchCustomState(channel.id).then((loaded) => {
            if (!active) {
                return;
            }
            setBaseline(loaded);
            setState(loaded);
            setLoading(false);
        }).catch(() => {
            // Don't seed an empty baseline on failure: a later save would clobber
            // stored data with blanks. Surface an error and block saving instead.
            if (active) {
                setLoadError(true);
                setLoading(false);
            }
        });

        return () => {
            active = false;
        };
    }, [channel.id]);

    const save = useCallback(async () => {
        if (loadError) {
            return;
        }
        const saved = await saveCustomState(channel.id, state);
        setBaseline(saved);
    }, [channel.id, state, loadError]);

    const reset = useCallback(() => {
        setState(baseline);
    }, [baseline]);

    // Keep the registered handlers pointing at the latest closures.
    const handlersRef = useRef({save, reset});
    handlersRef.current = {save, reset};

    useEffect(() => {
        registerHandlers({
            save: () => handlersRef.current.save(),
            reset: () => handlersRef.current.reset(),
        });
        return () => registerHandlers(null);
    }, [registerHandlers]);

    if (loading) {
        return (
            <div>
                <FormattedMessage
                    id='channel_settings.custom_tab.loading'
                    defaultMessage='Loading…'
                />
            </div>
        );
    }

    if (loadError) {
        return (
            <div style={{color: 'var(--error-text)'}}>
                <FormattedMessage
                    id='channel_settings.custom_tab.load_error'
                    defaultMessage='Could not load these settings. Reopen the tab to try again.'
                />
            </div>
        );
    }

    return (
        <div style={{color: 'var(--center-channel-color)'}}>
            <style>{NOTE_TEXTAREA_STYLES}</style>
            <h4 style={{color: 'var(--link-color)'}}>
                <FormattedMessage
                    id='channel_settings.custom_tab.heading'
                    defaultMessage='Advanced channel options'
                />
            </h4>
            <label htmlFor='demo-channel-note'>
                <FormattedMessage
                    id='channel_settings.custom_tab.note_label'
                    defaultMessage='Channel note'
                />
            </label>
            <textarea
                id='demo-channel-note'
                className={NOTE_TEXTAREA_CLASS}
                value={state.note}
                onChange={(e) => setState((prev) => ({...prev, note: e.target.value}))}
                placeholder={intl.formatMessage({id: 'channel_settings.custom_tab.note_placeholder', defaultMessage: 'Add a note for this channel…'})}
                rows={3}
                style={{display: 'block', width: '100%', marginTop: '4px'}}
            />
            <label style={{display: 'flex', alignItems: 'center', gap: '8px', marginTop: '12px'}}>
                <input
                    type='checkbox'
                    checked={state.pinGreeting}
                    onChange={(e) => setState((prev) => ({...prev, pinGreeting: e.target.checked}))}
                />
                <FormattedMessage
                    id='channel_settings.custom_tab.pin_greeting'
                    defaultMessage='Pin greeting message'
                />
            </label>
        </div>
    );
}

export default injectIntl(ChannelSettingsCustomTab);
