import {connect} from 'react-redux';

import {isReceiveWhatsappMessages} from 'selectors';

import LeftSidebarHeader from './left_sidebar_header';

const mapStateToProps = (state) => ({
    enabled: isReceiveWhatsappMessages(state),
});

export default connect(mapStateToProps)(LeftSidebarHeader);
