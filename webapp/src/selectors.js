import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isEnabled = (state) => getPluginState(state).enabled;

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;

export const subMenu = (state) => getPluginState(state).subMenu;

export const isReceiveWhatsappMessages = (state) =>
    getPluginState(state).preferences.whatsapp;

export const getActiveUsers = (state) => getPluginState(state).activeUsers.users;
