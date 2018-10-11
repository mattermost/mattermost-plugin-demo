import {combineReducers} from 'redux';

import {
    STATUS_CHANGE,
    OPEN_ROOT_MODAL,
    CLOSE_ROOT_MODAL,
    OPEN_SYSTEM_WIDE_SETTING_MODAL,
    CLOSE_SYSTEM_WIDE_SETTING_MODAL,
    SYSTEM_WIDE_SETTING_CHANGE,
} from './action_types';

const enabled = (state = false, action) => {
    switch (action.type) {
    case STATUS_CHANGE:
        return action.data;

    default:
        return state;
    }
};

const rootModalVisible = (state = false, action) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        return true;
    case CLOSE_ROOT_MODAL:
        return false;
    default:
        return state;
    }
};

const systemWideSettingModalVisible = (state = false, action) => {
    switch (action.type) {
        case OPEN_SYSTEM_WIDE_SETTING_MODAL:
            return true;
        case CLOSE_SYSTEM_WIDE_SETTING_MODAL:
            return false;
        default:
            return state;
    }
};

const systemWideSetting = (state = false, action) => {
    switch (action.type) {
        case SYSTEM_WIDE_SETTING_CHANGE:
            return action.data;
        default:
            return state;
    }
};

export default combineReducers({
    enabled,
    rootModalVisible,
    systemWideSettingModalVisible,
    systemWideSetting,
});

