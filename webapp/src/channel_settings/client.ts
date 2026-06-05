import {Client4} from 'mattermost-redux/client';

import type {Options} from '@mattermost/types/client4';

import {id as pluginId} from '../manifest';

import type {ChannelSettingsValues} from './types';

// Custom tab payload mirrors the server's customState struct.
export type CustomState = {
    note: string;
    pinGreeting: boolean;
};

function pluginRoute(): string {
    return `${Client4.getUrl()}/plugins/${pluginId}/channel_settings`;
}

// Build request options via Client4 so the CSRF token (X-CSRF-Token), credentials,
// and default headers are attached; without them cookie-authed POSTs are rejected
// (401) by the server proxy. Throw on non-2xx so the host keeps the tab dirty.
async function request<T>(url: string, options: Options): Promise<T> {
    const res = await fetch(url, Client4.getOptions(options));
    if (!res.ok) {
        throw new Error(`channel settings request failed: ${res.status}`);
    }
    return res.json() as Promise<T>;
}

export function fetchSchemaValues(channelId: string): Promise<ChannelSettingsValues> {
    return request<ChannelSettingsValues>(`${pluginRoute()}/${channelId}/schema`, {method: 'get'});
}

export function saveSchemaValues(channelId: string, values: ChannelSettingsValues): Promise<ChannelSettingsValues> {
    return request<ChannelSettingsValues>(`${pluginRoute()}/${channelId}/schema`, {
        method: 'post',
        body: JSON.stringify(values),
    });
}

export function fetchCustomState(channelId: string): Promise<CustomState> {
    return request<CustomState>(`${pluginRoute()}/${channelId}/custom`, {method: 'get'});
}

export function saveCustomState(channelId: string, state: CustomState): Promise<CustomState> {
    return request<CustomState>(`${pluginRoute()}/${channelId}/custom`, {
        method: 'post',
        body: JSON.stringify(state),
    });
}
