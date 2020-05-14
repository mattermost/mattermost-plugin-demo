import React from 'react';
import PropTypes from 'prop-types';

import {FormattedMessage} from 'react-intl';

export default class RHSView extends React.PureComponent {
    static propTypes = {
        team: PropTypes.object.isRequired,
    }

    render() {
        return (
            <div style={style.rhs}>
                <br/>
                <br/>
                <br/>
                <br/>
                <FormattedMessage
                    id='rhs.triggered'
                    defaultMessage='You have triggered the right-hand sidebar component of the demo plugin.'
                />
                <br/>
                <br/>
                <br/>
                <br/>
                <FormattedMessage
                    id='demo.testintl'
                    defaultMessage='This is the default string'
                />
                <br/>
                <br/>
                <br/>
                <br/>
                {'Links for custom routes'}
                <br/>
                <a onClick={() => window.WebappUtils.browserHistory.push('/plug/com.mattermost.demo-plugin/roottest')}>
                    {'/plug/com.mattermost.demo-plugin/roottest'}
                </a>
                <br/>
                <a onClick={() => window.WebappUtils.browserHistory.push(`/${this.props.team.name}/com.mattermost.demo-plugin/teamtest`)}>
                    {`/${this.props.team.name}/com.mattermost.demo-plugin/teamtest`}
                </a>
            </div>
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
    },
};
