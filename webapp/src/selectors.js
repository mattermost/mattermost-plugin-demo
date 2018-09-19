import Manifest from './manifest';

const getPluginState = (state) => state['plugins-' + Manifest.PluginId] || {};

export const isEnabled = (state) => getPluginState(state).enabled;

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;
