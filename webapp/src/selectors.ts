// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {id as pluginId} from './manifest';

type PluginState = {
    enabled?: boolean;
    rootModalVisible?: boolean;
    subMenu?: unknown;
};

type RootState = {
    [key: string]: PluginState | undefined;
};

const getPluginState = (state: RootState): PluginState => state['plugins-' + pluginId] || {};

export const isEnabled = (state: RootState) => Boolean(getPluginState(state).enabled);

export const isRootModalVisible = (state: RootState) => Boolean(getPluginState(state).rootModalVisible);

export const subMenu = (state: RootState) => getPluginState(state).subMenu;
