import React, {useCallback, useEffect, useState} from 'react';
import PropTypes from 'prop-types';

export default function ChannelSettingsSmokeTest({channel, setAreThereUnsavedChanges, registerSaveBarHandlers}) {
    const [value, setValue] = useState('');

    const handleSave = useCallback(async () => {
        // Smoke test: no server persistence; clearing dirty state matches a successful save.
        setAreThereUnsavedChanges?.(false);
    }, [setAreThereUnsavedChanges]);

    const handleReset = useCallback(() => {
        setValue('');
        setAreThereUnsavedChanges?.(false);
    }, [setAreThereUnsavedChanges]);

    /* eslint-disable consistent-return -- useEffect may return cleanup or nothing */
    useEffect(() => {
        if (!registerSaveBarHandlers) {
            return;
        }
        registerSaveBarHandlers({
            save: handleSave,
            reset: handleReset,
        });
        return () => registerSaveBarHandlers(null);
    }, [registerSaveBarHandlers, handleSave, handleReset]);
    /* eslint-enable consistent-return */

    const handleChange = useCallback((e) => {
        const newValue = e.target.value;
        setValue(newValue);
        setAreThereUnsavedChanges?.(newValue.length > 0);
    }, [setAreThereUnsavedChanges]);

    return (
        <div style={{padding: '20px'}}>
            <h3>{'Channel Settings Smoke Test'}</h3>
            <div style={{marginTop: '12px'}}>
                <strong>{'Display Name: '}</strong>{channel.display_name}
            </div>
            <div style={{marginTop: '4px'}}>
                <strong>{'Channel Name: '}</strong>{channel.name}
            </div>
            <div style={{marginTop: '4px'}}>
                <strong>{'Channel ID: '}</strong>{channel.id}
            </div>
            <div style={{marginTop: '16px'}}>
                <label htmlFor='smoke-test-input'>
                    <strong>{'Dirty-state test (type to mark dirty):'}</strong>
                </label>
                <br/>
                <input
                    id='smoke-test-input'
                    type='text'
                    value={value}
                    onChange={handleChange}
                    placeholder='Type here to mark tab as dirty'
                    style={{marginTop: '4px', padding: '6px', width: '300px'}}
                />
            </div>
        </div>
    );
}

ChannelSettingsSmokeTest.propTypes = {
    channel: PropTypes.object.isRequired,
    setAreThereUnsavedChanges: PropTypes.func,
    registerSaveBarHandlers: PropTypes.func,
};
