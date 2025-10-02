import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/common';
import {getMyPreferences} from 'mattermost-redux/selectors/entities/preferences';

import {id as PluginId} from './manifest';

import {
    CLOSE_ROOT_MODAL,
    OPEN_ROOT_MODAL,
    SET_ACTIVE_USERS,
    SET_WHATSAPP_PREF,
    STATUS_CHANGE,
    SUBMENU,
    SET_SESSION,
    SET_SESSION_ERROR,
    CLEAR_SESSION,
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

export const getBasePath = (state) => {
    const config = getConfig(state);

    let basePath = '/';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }
    return basePath;
};

export const getPluginServerRoute = (state) => {
    return getBasePath(state) + '/plugins/' + PluginId;
};

export const getStatus = () => async (dispatch, getState) => {
    fetch(getPluginServerRoute(getState()) + '/status').then((r) => r.json()).then((r) => {
        dispatch({
            type: STATUS_CHANGE,
            data: r.enabled,
        });
    });
};

export const getMyActiveSession = () => async (dispatch, getState) => {
    const state = getState();
    const url = `${getPluginServerRoute(state)}/sessions/active`;
    try {
        const res = await fetch(url, {method: 'GET'});
        if (res.status === 404) {
            dispatch({type: SET_SESSION, data: null});
            dispatch({type: SET_SESSION_ERROR, error: 'not_found'});
            return null;
        }
        if (!res.ok) {
            const text = await res.text();
            dispatch({type: SET_SESSION_ERROR, error: `http_${res.status}: ${text}`});
            return null;
        }
        const data = await res.json();
        dispatch({type: SET_SESSION, data});
        dispatch({type: SET_SESSION_ERROR, error: null});
        return data;
    } catch (err) {
        dispatch({type: SET_SESSION_ERROR, error: String(err)});
        return null;
    }
};

export const saveWhatsAppPreference = (enabled, options = {}) => async (dispatch, getState) => {
    const state = getState();
    const userId = getCurrentUserId(state);
    const url = `${getPluginServerRoute(state)}/whatsapp/preferences`;

    let enabledToSave = enabled;
    if (enabled === 'on') {
        const session = await dispatch(getMyActiveSession());
        const isActive = Boolean(session && !session.closedAt);
        if (!isActive) {
            // eslint-disable-next-line no-alert
            alert('No hay una sesión activa. La preferencia se mantendrá en Off.');
            enabledToSave = 'off';
            dispatch({type: SET_WHATSAPP_PREF, data: enabledToSave});
            return enabledToSave;
        }
    }

    try {
        const body = JSON.stringify({
            receive_notifications: enabledToSave === 'on',
            user_id: userId,
        });

        const response = await fetch(url, {
            method: 'PUT',
            headers: {'Content-Type': 'application/json'},
            body,
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP Error ${response.status}: ${errorText}`);
        }

        dispatch({type: SET_WHATSAPP_PREF, data: enabledToSave});
        return enabledToSave;
    } catch (error) {
        // eslint-disable-next-line no-console
        console.error('Error guardando la preferencia de WhatsApp:', error);
        return enabledToSave;
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

export const syncActiveUsers = (users) => async (dispatch) => {
    dispatch({
        type: SET_ACTIVE_USERS,
        data: users,
    });
};

export const websocketStatusChange = (message) => (dispatch) => dispatch({
    type: STATUS_CHANGE,
    data: message.data.enabled,
});
