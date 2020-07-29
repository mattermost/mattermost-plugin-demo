import React from 'react';
import {
  Switch,
  Route,
  useLocation
} from "react-router-dom";
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {useSelector} from 'react-redux';
import {id} from 'src/manifest';

export default function CustomSubRoutes() {
  const currentTeam = useSelector(getCurrentTeam);
  let query = useQuery();

  return (
    <Switch>
      <Route path={`/${currentTeam.name}/${id}/teamtest/subpath`}>
        <h3>Hello, Subpath Component!</h3>
      </Route>
      <Route path={`/${currentTeam.name}/${id}/teamtest/subpath-with-query`}>
        <h3>Hello, Subpath with Query!</h3>
        <p>The <code>search-value</code> in the query string is &quot;{query.get("search-value")}&quot;</p>
      </Route>
      <Route>
        <h3>Custom SubRoutes:</h3>
        <li>
          <a onClick={() => window.WebappUtils.browserHistory.push(`/${currentTeam.name}/${id}/teamtest/subpath`)}>Subpath</a>
        </li>
        <li>
          <a onClick={() => window.WebappUtils.browserHistory.push(`/${currentTeam.name}/${id}/teamtest/subpath-with-query?search-value=mattermost-plugin`)}>Subpath with Query</a>
        </li>
      </Route>
    </Switch>
  )
}

// A custom hook that builds on useLocation to parse the query string for you.
function useQuery() {
  return new URLSearchParams(useLocation().search);
}