// Local, structurally-identical mirror of the host's channel settings API.
// Plugins cannot import host internals, so these must stay in sync with
// webapp/channels/src/types/plugins/channel_settings.ts and
// webapp/channels/src/plugins/settings_schema/types.ts in the monorepo.

import type React from 'react';

import type {Channel} from '@mattermost/types/channels';

export type ChannelSettingsValues = {[name: string]: string};

export type RadioSettingOption = {
    value: string;
    text: string;
    helpText?: string;
};

export type RadioSetting = {
    name: string;
    type: 'radio';
    title?: string;
    helpText?: string;
    default: string;
    options: RadioSettingOption[];
};

// Custom setting components receive ONLY informChange; they must read the
// channel from Redux and hydrate themselves.
export type CustomSettingComponent = React.ComponentType<{informChange: (name: string, value: string) => void}>;

export type CustomSetting = {
    name: string;
    type: 'custom';
    title?: string;
    helpText?: string;
    default?: string;
    component: CustomSettingComponent;
};

export type Setting = RadioSetting | CustomSetting;

export type SettingsSection = {
    title: string;
    settings: Setting[];
    disabled?: boolean;
};

// Custom sections render with no props; read the channel from Redux if needed.
export type CustomSection = {
    title: string;
    component: React.ComponentType;
};

export type ChannelSettingsTabHandlers = {
    save: () => Promise<void>;
    reset: () => void;
};

// Host injects theme/webSocketClient beyond the documented contract.
export type ChannelSettingsTabBodyProps = {
    channel: Channel;
    setUnsaved: (unsaved: boolean) => void;
    registerHandlers: (handlers: ChannelSettingsTabHandlers | null) => void;
};

type ChannelSettingsTabBase = {
    uiName: string;
    icon?: string;
    shouldRender?: (state: unknown, channel: Channel) => boolean;
};

export type ChannelSettingsSchemaTab = ChannelSettingsTabBase & {
    sections: Array<SettingsSection | CustomSection>;
    onSave: (values: ChannelSettingsValues, channel: Channel) => Promise<void>;
    loadValues?: (channel: Channel) => ChannelSettingsValues | Promise<ChannelSettingsValues>;
    component?: never;
};

export type ChannelSettingsCustomTab = ChannelSettingsTabBase & {
    component: React.ComponentType<ChannelSettingsTabBodyProps>;
    sections?: never;
    onSave?: never;
};

export type ChannelSettingsTab = ChannelSettingsSchemaTab | ChannelSettingsCustomTab;
