import {connect} from 'react-redux';

import {GlobalState} from '@mattermost/types/store';
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';

import FilePreviewOverride from './file_preview_override';

const mapStateToProps = (state: GlobalState) => ({
    theme: getTheme(state),
});

export default connect(mapStateToProps)(FilePreviewOverride);
