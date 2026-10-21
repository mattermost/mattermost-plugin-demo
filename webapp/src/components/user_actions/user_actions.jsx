// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import {Button} from '@mattermost/compass-ui/components/button';

import '../demo_chrome.css';

export default class UserActionsComponent extends React.PureComponent {
    static propTypes = {
        openRootModal: PropTypes.func.isRequired,
        hide: PropTypes.func.isRequired,
    }

    onClick = () => {
        this.props.openRootModal();
        this.props.hide();
    }

    render() {
        return (
            <div className='DemoUserActions'>
                <FormattedMessage
                    id='useractions.demo'
                    defaultMessage='Demo Plugin: '
                />
                <Button
                    emphasis='primary'
                    size='small'
                    onClick={this.onClick}
                >
                    <FormattedMessage
                        id='useractions.action'
                        defaultMessage='Action'
                    />
                </Button>
            </div>
        );
    }
}
