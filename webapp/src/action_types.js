import {id as pluginId} from './manifest';

// Namespace your actions to avoid collisions.
export const STATUS_CHANGE = pluginId + '_status_change';

export const OPEN_ROOT_MODAL = pluginId + '_open_root_modal';
export const CLOSE_ROOT_MODAL = pluginId + '_close_root_modal';

export const SUBMENU = pluginId + '_submenu';

export const SET_WHATSAPP_PREF = pluginId + '_set_whats_pref';

export const SET_USER_PREFS = pluginId + '_set_user_pref';

export const SET_ACTIVE_USERS = pluginId + '_set_active_users';
