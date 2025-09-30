import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/common';
import {Client4} from 'mattermost-redux/client';

import {id as pluginId} from './manifest';
import {STATUS_CHANGE, OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL, SUBMENU, SET_WHATSAPP_PREF, SET_USER_PREFS} from './action_types';

export const openRootModal = (subMenuText = '') => (dispatch) => {
    dispatch({
        type: SUBMENU,
        subMenu: subMenuText,
    });
    dispatch({
        type: OPEN_ROOT_MODAL,
    });
};

export const closeRootModal = () => (dispatch) => {
    dispatch({
        type: CLOSE_ROOT_MODAL,
    });
};

export const mainMenuAction = openRootModal;
export const fileUploadMethodAction = openRootModal;
export const postDropdownMenuAction = openRootModal;
export const postDropdownSubMenuAction = openRootModal;
export const channelHeaderMenuAction = openRootModal;
export const fileDropdownMenuAction = openRootModal;

// TODO: Move this into mattermost-redux or mattermost-webapp.
export const getPluginServerRoute = (state) => {
    const config = getConfig(state);

    let basePath = '/';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + pluginId;
};

export const getStatus = () => async (dispatch, getState) => {
    fetch(getPluginServerRoute(getState()) + '/status').then((r) => r.json()).then((r) => {
        dispatch({
            type: STATUS_CHANGE,
            data: r.enabled,
        });
    });
};

export const saveWhatsAppPreference = (enabled) => async (dispatch, getState) => {
    const state = getState();
    const userId = getCurrentUserId(state);
    Client4.doFetch(`${getPluginServerRoute(state)}/whatsapp/preferences`, {
        method: 'put',
        body: JSON.stringify({
            receive_notifications: enabled === 'on',
            user_id: userId,
        }),
    });
    dispatch({
        type: SET_WHATSAPP_PREF,
        data: enabled,
    });
};

export const syncWhatsappPreferences = () => async (dispatch, getState) => {
    console.log('syncWhatsappPreferences');
    const state = getState();
    const userId = getCurrentUserId(state);
    Client4.doFetch(`${getPluginServerRoute(state)}/whatsapp/preferences`, {
        method: 'get',
        body: JSON.stringify({
            user_id: userId,
        }),
    }).then((r) => r.json()).then((r) => {
        console.log(r);
        dispatch({
            type: SET_USER_PREFS,
            data: r,
        });
    });
};

export const websocketStatusChange = (message) => (dispatch) => dispatch({
    type: STATUS_CHANGE,
    data: message.data.enabled,
});
