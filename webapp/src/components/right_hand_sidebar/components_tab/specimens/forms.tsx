// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';

import {Checkbox} from '@mattermost/compass-ui/components/checkbox';
import {Chip} from '@mattermost/compass-ui/components/chip';
import {Combobox} from '@mattermost/compass-ui/components/combobox';
import {DateRangePicker} from '@mattermost/compass-ui/components/date-range-picker';
import {Dropdown} from '@mattermost/compass-ui/components/dropdown';
import {Radio} from '@mattermost/compass-ui/components/radio';
import {SearchInput} from '@mattermost/compass-ui/components/search-input';
import {Select} from '@mattermost/compass-ui/components/select';
import {Switch} from '@mattermost/compass-ui/components/switch';
import {TextArea} from '@mattermost/compass-ui/components/text-area';
import {TextInput} from '@mattermost/compass-ui/components/text-input';

import {Row, Variant} from './shared';

const TEAM_OPTIONS = [
    {value: 'alpha', label: 'Alpha'},
    {value: 'bravo', label: 'Bravo'},
    {value: 'charlie', label: 'Charlie'},
];

export function SearchInputPreview() {
    return (
        <SearchInput label='Search'/>
    );
}

export function SearchInputDetail() {
    const [value, setValue] = useState('off-topic');
    return (
        <>
            <Variant label='Sizes'>
                <div className='DemoCompassShowcase__stack'>
                    <SearchInput
                        label='Small'
                        size='small'
                    />
                    <SearchInput
                        label='Medium'
                        size='medium'
                    />
                    <SearchInput
                        label='Large'
                        size='large'
                    />
                </div>
            </Variant>
            <Variant label='Filled / invalid'>
                <div className='DemoCompassShowcase__stack'>
                    <SearchInput
                        label='Search channels'
                        value={value}
                        onChange={(event) => setValue(event.target.value)}
                        onClear={() => setValue('')}
                    />
                    <SearchInput
                        invalid={true}
                        label='Invalid query'
                        value='???'
                    />
                </div>
            </Variant>
        </>
    );
}

export function TextInputPreview() {
    return (
        <TextInput label='Display name'/>
    );
}

export function TextInputDetail() {
    return (
        <Variant label='Default / invalid'>
            <div className='DemoCompassShowcase__stack'>
                <TextInput label='Display name'/>
                <TextInput
                    invalid={true}
                    label='Channel URL'
                    defaultValue='not valid'
                />
            </div>
        </Variant>
    );
}

export function TextAreaPreview() {
    return (
        <TextArea label='Purpose'/>
    );
}

export function TextAreaDetail() {
    return (
        <TextArea
            defaultValue='A place to discuss Compass UI.'
            label='Channel purpose'
            maxLength={120}
            showCharacterCount={true}
        />
    );
}

export function CheckboxPreview() {
    return (
        <Checkbox defaultChecked={true}>{'Notify me'}</Checkbox>
    );
}

export function CheckboxDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <Checkbox>{'Unchecked'}</Checkbox>
            <Checkbox defaultChecked={true}>{'Checked'}</Checkbox>
            <Checkbox indeterminate={true}>{'Indeterminate'}</Checkbox>
            <Checkbox
                disabled={true}
                defaultChecked={true}
            >
                {'Disabled'}
            </Checkbox>
        </div>
    );
}

export function RadioPreview() {
    return (
        <Radio
            defaultChecked={true}
            name='preview-plan'
        >
            {'Professional'}
        </Radio>
    );
}

export function RadioDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <Radio
                defaultChecked={true}
                name='plan'
                value='free'
            >
                {'Free'}
            </Radio>
            <Radio
                name='plan'
                value='professional'
            >
                {'Professional'}
            </Radio>
            <Radio
                disabled={true}
                name='plan'
                value='enterprise'
            >
                {'Enterprise'}
            </Radio>
        </div>
    );
}

export function SwitchPreview() {
    return (
        <Switch defaultChecked={true}>{'Desktop notifications'}</Switch>
    );
}

export function SwitchDetail() {
    return (
        <div className='DemoCompassShowcase__stack'>
            <Switch>{'Off'}</Switch>
            <Switch defaultChecked={true}>{'On'}</Switch>
            <Switch
                disabled={true}
                defaultChecked={true}
            >
                {'Disabled'}
            </Switch>
        </div>
    );
}

export function SelectPreview() {
    return (
        <Select
            defaultValue='alpha'
            label='Team'
            options={TEAM_OPTIONS}
        />
    );
}

export function SelectDetail() {
    const [value, setValue] = useState('alpha');
    return (
        <Select
            label='Team'
            options={TEAM_OPTIONS}
            value={value}
            onChange={setValue}
        />
    );
}

export function ChipPreview() {
    return (
        <Chip>{'UX Design'}</Chip>
    );
}

export function ChipDetail() {
    return (
        <Row>
            <Chip>{'Default'}</Chip>
            <Chip colored={true}>{'Colored'}</Chip>
            <Chip onRemove={() => undefined}>{'Removable'}</Chip>
            <Chip error={true}>{'Error'}</Chip>
        </Row>
    );
}

export function ComboboxPreview() {
    return (
        <Combobox
            label='Members'
            options={TEAM_OPTIONS}
            placeholder='Find a member'
        />
    );
}

export function ComboboxDetail() {
    return (
        <Combobox
            label='Members'
            multiple={true}
            options={TEAM_OPTIONS}
            placeholder='Add people'
        />
    );
}

export function DropdownPreview() {
    return (
        <Dropdown>{'Town Square'}</Dropdown>
    );
}

export function DropdownDetail() {
    return (
        <Variant label='Closed / open'>
            <Row>
                <Dropdown>{'Town Square'}</Dropdown>
                <Dropdown isOpen={true}>{'Off-Topic'}</Dropdown>
            </Row>
        </Variant>
    );
}

export function DateRangePickerPreview() {
    return (
        <DateRangePicker placeholder='Select a date'/>
    );
}

export function DateRangePickerDetail() {
    return (
        <Variant label='Date / range'>
            <div className='DemoCompassShowcase__stack'>
                <DateRangePicker
                    placeholder='Date'
                    value='2026-09-02'
                />
                <DateRangePicker
                    mode='range'
                    placeholder='Date range'
                />
            </div>
        </Variant>
    );
}
