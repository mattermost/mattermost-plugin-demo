import {connect} from 'react-redux';

import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {getUnreadChannels} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentUserId, getMyChannelMemberships} from 'mattermost-redux/selectors/entities/common';

import {getActiveUsers, isEnabled} from 'selectors';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
    team: getCurrentTeam(state),
    unreadChannels: getUnreadChannels(state),
    myChannelMemberships: getMyChannelMemberships(state),
    currentUserId: getCurrentUserId(state),
    activeUsers: getActiveUsers(state),
    reduxState: state,
});

export default connect(mapStateToProps)(RHSView);
