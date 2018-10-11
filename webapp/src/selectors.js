import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isEnabled = (state) => getPluginState(state).enabled;

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;

export const isSystemWideSettingModalVisible = (state) => getPluginState(state).systemWideSettingModalVisible;
export const getSetting = (state) => getPluginState(state).systemWideSetting;
