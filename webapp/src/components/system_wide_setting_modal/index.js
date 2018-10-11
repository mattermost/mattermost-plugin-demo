import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {closeSystemWideSettingModal} from 'actions';
import {isSystemWideSettingModalVisible, getSetting} from 'selectors';

import SystemWideSettingModal from './system_wide_setting_modal';

const mapStateToProps = (state) => ({
    visible: isSystemWideSettingModalVisible(state),
    setting: getSetting(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeSystemWideSettingModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(SystemWideSettingModal);
