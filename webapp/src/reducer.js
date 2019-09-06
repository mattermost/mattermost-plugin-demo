import {combineReducers} from 'redux';

import {STATUS_CHANGE, OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL, SUBMENU} from './action_types';

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

export default combineReducers({
    enabled,
    rootModalVisible,
    subMenu,
});

