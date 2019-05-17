import {connect} from 'react-redux';

import {isEnabled} from 'selectors';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
});

export default connect(mapStateToProps)(RHSView);
