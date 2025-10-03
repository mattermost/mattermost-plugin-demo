import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/common';
import {getMyPreferences} from 'mattermost-redux/selectors/entities/preferences';
import {getCurrentTeamUrl} from 'mattermost-redux/selectors/entities/teams';

import {id as PluginId} from './manifest';

import {
    CLOSE_ROOT_MODAL,
    OPEN_ROOT_MODAL,
    SET_ACTIVE_USERS,
    SET_WHATSAPP_PREF,
    STATUS_CHANGE,
    SUBMENU,
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

async function getSessionIdByUserId(baseUrl, userId) {
    
    const res = await fetch(`${baseUrl}/sessions/${encodeURIComponent(userId)}`, {method: 'GET'});
    if (!res.ok) {
        return null;
    }
    const data = await res.json();
    return data?.id || null;
}

export const saveWhatsAppPreference = (enabled) => async (dispatch, getState) => {
    const state = getState();
    const userId = getCurrentUserId(state);
    const baseUrl = getPluginServerRoute(state);
    const isOn = enabled === 'on';
    try {
        if (isOn) {
            const res = await fetch(`${baseUrl}/sessions`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({userID: userId}),
            });
            if (!res.ok) {
                const errText = await res.text().catch(() => '');
                throw new Error(`Error creando sesión: ${errText || res.statusText}`);
            }
        } else {
            // Close by userID directly per server API
            const delRes = await fetch(`${baseUrl}/sessions/${encodeURIComponent(userId)}`, {method: 'DELETE'});
            if (!delRes.ok && delRes.status !== 404) {
                const errText = await delRes.text().catch(() => '');
                throw new Error(`Error eliminando sesión: ${errText || delRes.statusText}`);
            }
        }
        dispatch({
            type: SET_WHATSAPP_PREF,
            data: enabled,
        });
    } catch (error) {
        // eslint-disable-next-line no-console
        console.error('Error guardando preferencia WhatsApp:', error);
    }
};

// cuando elimino por nerlo en null  si el objero session exite deber ser en on

export const syncWhatsappPreferences = () => async (dispatch, getState) => {
    const state = getState();
    const baseUrl = getPluginServerRoute(state);
    const userId = getCurrentUserId(state);

    try {
        const res = await fetch(`${baseUrl}/sessions/${encodeURIComponent(userId)}`, {method: 'GET'});
        if (res.ok) {
            dispatch({
                type: SET_WHATSAPP_PREF,
                data: 'on',
            });
            return;
        }
        if (res.status === 404) {
            dispatch({
                type: SET_WHATSAPP_PREF,
                data: 'off',
            });
            return;
        }
    } catch (e) {
        console.error('Error sincronizando preferencia WhatsApp:', e);
    }

    
    const PreSavedPreferenceList = getMyPreferences(state);
    const whatsappSetting = PreSavedPreferenceList[`pp_${PluginId}--${PREFERENCE_NAME_WHATSAPP}`];
    if (whatsappSetting?.value) {
        dispatch({
            type: SET_WHATSAPP_PREF,
            data: whatsappSetting.value,
        });
    }
};// caniar y crear un fech al servido para preguntar si el usuario tiene una session

export const getActiveUsers = () => async (dispatch, getState) => {
    const state = getState();
    const url = `${getPluginServerRoute(state)}/sessions`;

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
