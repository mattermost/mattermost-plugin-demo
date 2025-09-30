import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/common';
import {getMyPreferences} from 'mattermost-redux/selectors/entities/preferences';

import {id as PluginId} from './manifest';

import {
    STATUS_CHANGE,
    OPEN_ROOT_MODAL,
    CLOSE_ROOT_MODAL,
    SUBMENU,
    SET_WHATSAPP_PREF,
    SET_ACTIVE_USERS,
} from './action_types';
import {PREFERENCE_NAME_WHATSAPP} from './constants';

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

    return basePath + '/plugins/' + PluginId;
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
    const url = `${getPluginServerRoute(state)}/whatsapp/preferences`;

    try {
        const body = JSON.stringify({
            receive_notifications: enabled === 'on',
            user_id: userId,
        });

        const response = await fetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body,
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP Error ${response.status}: ${errorText}`);
        }

        dispatch({
            type: SET_WHATSAPP_PREF,
            data: enabled,
        });
    } catch (error) {
        console.error('Error guardando la preferencia de WhatsApp:', error);
    }
};

export const syncWhatsappPreferences = () => async (dispatch, getState) => {
    const state = getState();

    const PreSavedPreferenceList = getMyPreferences(state);

    const whatsappSetting = PreSavedPreferenceList[`pp_${PluginId}--${PREFERENCE_NAME_WHATSAPP}`];

    dispatch({
        type: SET_WHATSAPP_PREF,
        data: whatsappSetting.value,
    });
};

export const getActiveUsers = () => async (dispatch, getState) => {
    const state = getState();
    const url = `${getPluginServerRoute(state)}/whatsapp/enabled/users`;

    return fetch(url, {method: 'GET'}).
        then((r) => r.json()).
        then((data) => {
            dispatch({
                type: SET_ACTIVE_USERS,
                data: data.active_users,
            });
        }).
        catch((error) => {
            console.error('Error obteniendo los usuarios:', error);
        });
};

export const syncActiveUsers = (users) => async (dispatch, _) => {
    dispatch({
        type: SET_ACTIVE_USERS,
        data: users,
    });
};

export const websocketStatusChange = (message) => (dispatch) => dispatch({
    type: STATUS_CHANGE,
    data: message.data.enabled,
});
