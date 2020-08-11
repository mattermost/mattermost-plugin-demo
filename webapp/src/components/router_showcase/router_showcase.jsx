import React from 'react';
import {
    Switch,
    Route,
    useLocation,
} from 'react-router-dom';
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {useSelector} from 'react-redux';

import {id} from 'src/manifest';

export default function RouterShowcase() {
    const currentTeam = useSelector(getCurrentTeam);
    const query = useQuery();

    function pushTeamUrlPath(path) {
        return window.WebappUtils.browserHistory.push(teamUrlPath(path));
    }

    function teamUrlPath(path) {
        return `/${currentTeam.name}/${id}/teamtest/${path}`;
    }

    return (
        <Switch>
            <Route path={teamUrlPath('subpath')}>
                <h3>{'Hello, Subpath Component!'}</h3>
            </Route>
            <Route path={teamUrlPath('subpath-with-query')}>
                <h3>{'Hello, Subpath with Query!'}</h3>
                <p>{`The search-value in the query string is ${query.get('search-value') === 'on' ? 'enabled' : 'disabled'}.`}</p>
            </Route>
            <Route>
                <h3>{'Custom SubRoutes:'}</h3>
                <li>
                    <a onClick={() => pushTeamUrlPath('subpath')}>{'Subpath'}</a>
                </li>
                <li>
                    <a onClick={() => pushTeamUrlPath('subpath-with-query?search-value=on')}>{'Subpath with Query'}</a>
                </li>
            </Route>
        </Switch>
    );
}

// A custom hook that builds on useLocation to parse the query string for you
function useQuery() {
    return new URLSearchParams(useLocation().search);
}
