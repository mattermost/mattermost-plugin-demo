import {combineReducers} from 'redux';

import {
    STATUS_CHANGE,
    OPEN_ROOT_MODAL,
    CLOSE_ROOT_MODAL,
    SUBMENU,
    SET_WHATSAPP_PREF,
    SET_USER_PREFS,
    SET_ACTIVE_USERS,
} from './action_types';

const preferencesInitialState = {
    whatsapp: false,
};

const activeUsersInitialState = {
    users: [],
};

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

const subMenu = (state = '', action) => {
    switch (action.type) {
    case SUBMENU:
        return action.subMenu;

    default:
        return state;
    }
};

const preferences = (state = preferencesInitialState, action) => {
    switch (action.type) {
    case SET_WHATSAPP_PREF:
        return {
            ...state,
            whatsapp: action.data === 'on',
        };
    case SET_USER_PREFS:
        return {
            ...state,
            ...action.data,
        };
    default:
        return state;
    }
};

const activeUsers = (state = activeUsersInitialState, action) => {
    switch (action.type) {
    case SET_ACTIVE_USERS:
        return {
            ...state,
            users: action.data,
        };
    default:
        return state;
    }
};

export default combineReducers({
    enabled,
    rootModalVisible,
    subMenu,
    preferences,
    activeUsers,
});

