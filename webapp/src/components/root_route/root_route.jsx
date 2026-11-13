import React from 'react';
import {FormattedMessage, useIntl} from 'react-intl';
import {useHistory, useLocation} from 'react-router-dom';

import {Button} from '@mattermost/compass-ui/components/button';

import PluginPage, {Facts} from '../plugin_page/plugin_page';

export default function RootRoute() {
    const {formatMessage} = useIntl();
    const history = useHistory();
    const location = useLocation();

    return (
        <PluginPage
            title={formatMessage({id: 'rootRoute.title', defaultMessage: 'Root route'})}
            description={formatMessage(
                {
                    id: 'rootRoute.description',
                    defaultMessage: 'Registered with registerCustomRoute. It lives at {path}, outside any team, so there is no team sidebar or current channel. Use it for standalone pages like onboarding or OAuth callbacks.',
                },
                {path: '/plug/<plugin id>/<route>'},
            )}
        >
            <Facts items={[{label: formatMessage({id: 'route.currentPath', defaultMessage: 'Current path'}), value: location.pathname}]}/>
            <div className='demo-page__row'>
                <Button
                    emphasis='secondary'
                    onClick={() => history.push('/')}
                >
                    <FormattedMessage
                        id='rootRoute.back'
                        defaultMessage='Back to Mattermost'
                    />
                </Button>
            </div>
        </PluginPage>
    );
}
