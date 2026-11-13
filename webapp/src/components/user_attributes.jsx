// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {Button} from '@mattermost/compass-ui/components/button';

export default class UserAttributes extends React.PureComponent {
    static propTypes = {
        hide: PropTypes.func.isRequired,
    }

    onClick = () => {
        this.props.hide();
    }

    render() {
        return (
            <div className='DemoUserAttributes'>
                <Button
                    emphasis='link'
                    size='small'
                    onClick={this.onClick}
                >
                    {'Demo Plugin: User Attributes'}
                </Button>
            </div>
        );
    }
}
