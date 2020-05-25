import {connect} from 'react-redux';

import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';

import {isEnabled} from 'selectors';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
    team: getCurrentTeam(state),
});

export default connect(mapStateToProps)(RHSView);
