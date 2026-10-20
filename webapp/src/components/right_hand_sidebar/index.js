import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';

import {openRootModal} from 'actions';
import {isEnabled} from 'selectors';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
    team: getCurrentTeam(state),
    channel: getCurrentChannel(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({openRootModal}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(RHSView);
