// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import {PlugIcon} from '../icons';

import '../demo_chrome.css';

// LeftSidebarHeader is a pure component, later connected to the Redux store so as to
// show the plugin's enabled / disabled status.
export default class LeftSidebarHeader extends React.PureComponent {
    static propTypes = {
        enabled: PropTypes.bool.isRequired,
    }

    render() {
        return (
            <div
                className='DemoLeftSidebarHeader'
                data-testid='demo-plugin-hook-status'
            >
                <span className='DemoLeftSidebarHeader__icon'>
                    <PlugIcon size='16'/>
                </span>
                <FormattedMessage
                    id='sidebar.demo'
                    defaultMessage='Demo Plugin:'
                />
                {' '}
                {this.props.enabled ?
                    <span>
                        <FormattedMessage
                            id='sidebar.enabled'
                            defaultMessage='Enabled'
                        />
                    </span> :
                    <span>
                        <FormattedMessage
                            id='sidebar.disabled'
                            defaultMessage='Disabled'
                        />
                    </span>
                }
            </div>
        );
    }
}
