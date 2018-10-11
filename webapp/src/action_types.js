import {id as pluginId} from './manifest';

// Namespace your actions to avoid collisions.
export const STATUS_CHANGE = pluginId + '_status_change';

export const OPEN_ROOT_MODAL = pluginId + '_open_root_modal';
export const CLOSE_ROOT_MODAL = pluginId + '_close_root_modal';

export const OPEN_SYSTEM_WIDE_SETTING_MODAL = pluginId + '_open_system_wide_setting_modal';
export const CLOSE_SYSTEM_WIDE_SETTING_MODAL = pluginId + '_close_system_wide_setting_modal';
export const SYSTEM_WIDE_SETTING_CHANGE = pluginId + '_system_wide_setting_change';
