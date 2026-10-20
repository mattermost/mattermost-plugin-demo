import React from 'react';
import {FormattedMessage, defineMessages, useIntl} from 'react-intl';
import {
    Switch,
    Route,
    useHistory,
    useLocation,
    useParams,
    useRouteMatch,
} from 'react-router-dom';

import {Switch as Toggle} from '@mattermost/compass-ui/components/switch';
import {Tabs} from '@mattermost/compass-ui/components/tabs';
import {Tag} from '@mattermost/compass-ui/components/tag';

import PluginPage, {Facts} from '../plugin_page/plugin_page';

const messages = defineMessages({
    title: {id: 'teamRoute.title', defaultMessage: 'Team route'},
    description: {id: 'teamRoute.description', defaultMessage: 'Registered with registerNeedsTeamRoute. It renders in the center of a team, so team context is available, and it can own nested routes below its base path.'},
    overview: {id: 'teamRoute.tab.overview', defaultMessage: 'Overview'},
    subpath: {id: 'teamRoute.tab.subpath', defaultMessage: 'Subpath'},
    query: {id: 'teamRoute.tab.query', defaultMessage: 'Query string'},
    on: {id: 'teamRoute.query.on', defaultMessage: 'on'},
    off: {id: 'teamRoute.query.off', defaultMessage: 'off'},
    empty: {id: 'teamRoute.query.empty', defaultMessage: '(empty)'},
    team: {id: 'teamRoute.fact.team', defaultMessage: 'Team'},
    basePath: {id: 'teamRoute.fact.basePath', defaultMessage: 'Base path'},
    matchedRoute: {id: 'teamRoute.fact.matchedRoute', defaultMessage: 'Matched route'},
    currentPath: {id: 'route.currentPath', defaultMessage: 'Current path'},
});

const SUBPAGES = [
    {key: 'overview', path: ''},
    {key: 'subpath', path: '/subpath'},
    {key: 'query', path: '/subpath-with-query'},
];

export default function RouterShowcase() {
    const {formatMessage} = useIntl();
    const {team} = useParams();
    const {url} = useRouteMatch();
    const history = useHistory();
    const location = useLocation();

    const base = url.replace(/\/$/, '');
    const active = [...SUBPAGES].reverse().find((page) => location.pathname.startsWith(base + page.path));
    const searchEnabled = new URLSearchParams(location.search).get('search-value') === 'on';

    const tabs = SUBPAGES.map((page) => ({
        key: page.key,
        label: formatMessage(messages[page.key]),
        id: `demo-route-tab-${page.key}`,
        panelId: 'demo-route-panel',
    }));

    return (
        <PluginPage
            title={formatMessage(messages.title)}
            description={formatMessage(messages.description)}
            nav={(
                <Tabs
                    tabs={tabs}
                    activeKey={active.key}
                    onChange={(key) => history.push(base + SUBPAGES.find((page) => page.key === key).path)}
                    appearance='underlined'
                    className='demo-page__tabs'
                />
            )}
        >
            <div
                id='demo-route-panel'
                role='tabpanel'
                aria-labelledby={`demo-route-tab-${active.key}`}
                className='demo-page__section'
            >
                <Switch>
                    <Route path={`${base}/subpath-with-query`}>
                        <p className='demo-page__text'>
                            <FormattedMessage
                                id='teamRoute.query.text'
                                defaultMessage='Query parameters survive navigation. Toggle the switch to rewrite the URL with history.replace and read it back with useLocation.'
                            />
                        </p>
                        <div className='demo-page__row'>
                            <Toggle
                                checked={searchEnabled}
                                onChange={(e) => history.replace(`${base}/subpath-with-query${e.target.checked ? '?search-value=on' : ''}`)}
                            >
                                {'search-value'}
                            </Toggle>
                            <Tag
                                label={formatMessage(searchEnabled ? messages.on : messages.off)}
                                type={searchEnabled ? 'success' : 'default'}
                                size='small'
                            />
                        </div>
                        <Facts items={[{label: formatMessage(messages.query), value: location.search || formatMessage(messages.empty)}]}/>
                    </Route>
                    <Route path={`${base}/subpath`}>
                        <p className='demo-page__text'>
                            <FormattedMessage
                                id='teamRoute.subpath.text'
                                defaultMessage='This view comes from a nested {route} inside the plugin component. The host only matches the base path; everything after it belongs to the plugin.'
                                values={{route: <code key='route'>{'<Route>'}</code>}}
                            />
                        </p>
                        <Facts items={[{label: formatMessage(messages.matchedRoute), value: `${base}/subpath`}]}/>
                    </Route>
                    <Route>
                        <Facts
                            items={[
                                {label: formatMessage(messages.team), value: team},
                                {label: formatMessage(messages.basePath), value: base},
                                {label: formatMessage(messages.currentPath), value: location.pathname},
                            ]}
                        />
                    </Route>
                </Switch>
            </div>
        </PluginPage>
    );
}
