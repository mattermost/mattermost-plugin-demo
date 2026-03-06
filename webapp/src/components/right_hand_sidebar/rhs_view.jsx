import React, {useEffect, useState} from 'react';
import PropTypes from 'prop-types';

import {FormattedMessage} from 'react-intl';

import {id as pluginId} from '../../manifest';

export default function RHSView({team, channel}) {
    const [autoPopout, setAutoPopout] = useState(false);
    const popoutSupported = Boolean(window.WebappUtils?.popouts?.popoutRhsPlugin);

    // useEffect example: automatically pop out the RHS when autoPopout is toggled on
    useEffect(() => {
        if (autoPopout) {
            if (popoutSupported) {
                window.WebappUtils.popouts.popoutRhsPlugin('Demo Plugin', pluginId, team.name, channel.name);
            }
            setAutoPopout(false);
        }
    }, [autoPopout, popoutSupported, team.name, channel.name]);

    const handlePopout = () => {
        if (popoutSupported) {
            window.WebappUtils.popouts.popoutRhsPlugin('Demo Plugin', pluginId, team.name, channel.name);
        }
    };

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
            <a onClick={() => window.WebappUtils.browserHistory.push(`/${team.name}/com.mattermost.demo-plugin/teamtest`)}>
                {`/${team.name}/com.mattermost.demo-plugin/teamtest`}
            </a>
            <br/>
            <br/>
            <hr/>
            <br/>
            <strong>{'Pop Out RHS Demo'}</strong>
            <br/>
            <br/>
            <button
                className='btn btn-primary'
                disabled={!popoutSupported}
                onClick={handlePopout}
            >
                <FormattedMessage
                    id='rhs.popout'
                    defaultMessage='Pop Out RHS'
                />
            </button>
            <br/>
            <br/>
            <button
                className='btn btn-tertiary'
                disabled={!popoutSupported}
                onClick={() => setAutoPopout(true)}
            >
                <FormattedMessage
                    id='rhs.popout.useeffect'
                    defaultMessage='Pop Out via useEffect'
                />
            </button>
        </div>
    );
}

RHSView.propTypes = {
    team: PropTypes.object.isRequired,
    channel: PropTypes.object.isRequired,
};

const style = {
    rhs: {
        padding: '10px',
    },
};
