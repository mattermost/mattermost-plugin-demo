import React, {useState} from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage, defineMessages, useIntl} from 'react-intl';

import {Button} from '@mattermost/compass-ui/components/button';
import {Divider} from '@mattermost/compass-ui/components/divider';
import {SearchInput} from '@mattermost/compass-ui/components/search-input';
import OpenInNewIcon from '@mattermost/compass-icons/components/open-in-new';
import PlayIcon from '@mattermost/compass-icons/components/play';

import {id as pluginId} from 'manifest';

import {matchesQuery, sections} from './surfaces';

const messages = defineMessages({
    showMe: {id: 'rhs.guide.action.show', defaultMessage: 'Show me'},
    popOut: {id: 'rhs.guide.action.popout', defaultMessage: 'Pop out'},
    open: {id: 'rhs.guide.action.open', defaultMessage: 'Open'},
    filterPlaceholder: {id: 'rhs.guide.filter.placeholder', defaultMessage: 'Filter by name or API'},
    filterLabel: {id: 'rhs.guide.filter.label', defaultMessage: 'Filter plugin surfaces'},
});

const navigate = (path) => window.WebappUtils.browserHistory.push(path);

function RouteLink({path}) {
    return (
        <a
            className='demo-rhs__path'
            href={path}
            onClick={(e) => {
                e.preventDefault();
                navigate(path);
            }}
        >
            {path}
        </a>
    );
}

RouteLink.propTypes = {
    path: PropTypes.string.isRequired,
};

export default function Guide({team, channel, openRootModal}) {
    const {formatMessage} = useIntl();
    const [query, setQuery] = useState('');

    const rootRoute = `/plug/${pluginId}/roottest`;

    const popouts = window.WebappUtils?.popouts;
    const canPopout = Boolean(popouts?.popoutRhsPlugin && team && channel) && popouts.canPopout() && !popouts.isPopoutWindow();

    const open = formatMessage(messages.open);
    const route = (path) => ({label: open, icon: <OpenInNewIcon/>, run: () => navigate(path), path});
    const actions = {
        openRoot: {label: formatMessage(messages.showMe), icon: <PlayIcon/>, run: () => openRootModal()},
        popout: {
            label: formatMessage(messages.popOut),
            icon: <OpenInNewIcon/>,
            disabled: !canPopout,
            run: () => popouts.popoutRhsPlugin(formatMessage({id: 'plugin.name', defaultMessage: 'Demo Plugin'}), pluginId, team.name, channel.name),
        },
        openAdminConsole: {label: open, icon: <OpenInNewIcon/>, run: () => navigate(`/admin_console/plugins/plugin_${pluginId}`)},
        openRootRoute: route(rootRoute),
        ...(team && {
            openTeamRoute: route(`/${team.name}/${pluginId}/teamtest`),
            openTicketsRoute: route(`/${team.name}/${pluginId}/tickets`),
        }),
    };

    const visibleSections = sections.
        map((section) => ({...section, items: section.items.filter((item) => matchesQuery(item, query))})).
        filter((section) => section.items.length > 0);

    return (
        <div className='demo-rhs__guide'>
            <header className='demo-rhs__header'>
                <p className='demo-rhs__text'>
                    <FormattedMessage
                        id='rhs.guide.intro'
                        defaultMessage='Below is a list of everything this plugin adds to Mattermost, and where to find it.'
                    />
                </p>
                <SearchInput
                    size='small'
                    placeholder={formatMessage(messages.filterPlaceholder)}
                    aria-label={formatMessage(messages.filterLabel)}
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onClear={() => setQuery('')}
                />
            </header>

            {visibleSections.length === 0 && (
                <p className='demo-rhs__text demo-rhs__empty'>
                    <FormattedMessage
                        id='rhs.guide.empty'
                        defaultMessage='No surfaces match "{query}".'
                        values={{query}}
                    />
                </p>
            )}

            {visibleSections.map((section, index) => (
                <section
                    key={section.id}
                    className='demo-rhs__section'
                    aria-labelledby={`demo-rhs-${section.id}`}
                >
                    {index > 0 && <Divider/>}
                    <h3
                        id={`demo-rhs-${section.id}`}
                        className='demo-rhs__section-title'
                    >
                        {section.title}
                    </h3>
                    <ul className='demo-rhs__list'>
                        {section.items.map((item) => {
                            const action = actions[item.action];
                            return (
                                <li
                                    key={item.name}
                                    className='demo-rhs__item'
                                >
                                    <div className='demo-rhs__item-text'>
                                        <span className='demo-rhs__name'>{item.name}</span>
                                        <code className='demo-rhs__api'>{item.api}</code>
                                        <span className='demo-rhs__text'>{item.where}</span>
                                        {action?.path && <RouteLink path={action.path}/>}
                                        {item.action === 'showTranslation' && (
                                            <span className='demo-rhs__sample'>
                                                <FormattedMessage
                                                    id='demo.testintl'
                                                    defaultMessage='This is the default string'
                                                />
                                            </span>
                                        )}
                                    </div>
                                    {action && (
                                        <Button
                                            size='x-small'
                                            emphasis='tertiary'
                                            leadingIcon={action.icon}
                                            disabled={action.disabled}
                                            onClick={action.run}
                                        >
                                            {action.label}
                                        </Button>
                                    )}
                                </li>
                            );
                        })}
                    </ul>
                </section>
            ))}
        </div>
    );
}

Guide.propTypes = {
    team: PropTypes.object,
    channel: PropTypes.object,
    openRootModal: PropTypes.func.isRequired,
};
