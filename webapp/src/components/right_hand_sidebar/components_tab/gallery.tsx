import React, {useMemo, useState} from 'react';
import {defineMessages, useIntl} from 'react-intl';

import {EmptyState} from '@mattermost/compass-ui/components/empty-state';
import {Scrollbar} from '@mattermost/compass-ui/components/scrollbar';
import {SearchInput} from '@mattermost/compass-ui/components/search-input';

import {CATALOG} from './catalog';
import {filterCatalog} from './helpers';
import ScaledPreview from './scaled_preview';
import {CATEGORIES, CATEGORY_LABELS, type CatalogEntry} from './types';

const messages = defineMessages({
    search: {id: 'components.search', defaultMessage: 'Search components'},
    emptyTitle: {id: 'components.empty.title', defaultMessage: 'No components match'},
    emptyDescription: {id: 'components.empty.description', defaultMessage: 'Try a different search.'},
    clearSearch: {id: 'components.empty.clear', defaultMessage: 'Clear search'},
});

type Props = {
    onSelect: (entry: CatalogEntry, trigger: HTMLButtonElement) => void;
};

export default function Gallery({onSelect}: Props) {
    const {formatMessage} = useIntl();
    const [query, setQuery] = useState('');
    const entries = useMemo(() => filterCatalog(CATALOG, query), [query]);

    return (
        <div className='DemoCompassShowcase'>
            <div className='DemoCompassShowcase__chrome'>
                <SearchInput
                    size='small'
                    placeholder={formatMessage(messages.search)}
                    aria-label={formatMessage(messages.search)}
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                    onClear={() => setQuery('')}
                />
            </div>
            <Scrollbar className='DemoCompassShowcase__scroll'>
                {entries.length === 0 ? (
                    <EmptyState
                        title={formatMessage(messages.emptyTitle)}
                        description={formatMessage(messages.emptyDescription)}
                        action={{children: formatMessage(messages.clearSearch), onClick: () => setQuery('')}}
                    />
                ) : (
                    <div className='DemoCompassShowcase__sections'>
                        {CATEGORIES.map((category) => {
                            const items = entries.filter((entry) => entry.category === category);
                            if (items.length === 0) {
                                return null;
                            }
                            const headingId = `compass-showcase-${category}`;
                            return (
                                <section
                                    key={category}
                                    className='DemoCompassShowcase__section'
                                    aria-labelledby={headingId}
                                >
                                    <h2
                                        id={headingId}
                                        className='DemoCompassShowcase__sectionTitle'
                                    >
                                        {CATEGORY_LABELS[category]}
                                    </h2>
                                    <ul className='DemoCompassShowcase__grid'>
                                        {items.map((entry) => (
                                            <CatalogCard
                                                key={entry.id}
                                                entry={entry}
                                                onSelect={onSelect}
                                            />
                                        ))}
                                    </ul>
                                </section>
                            );
                        })}
                    </div>
                )}
            </Scrollbar>
        </div>
    );
}

type CardProps = {
    entry: CatalogEntry;
    onSelect: Props['onSelect'];
};

function CatalogCard({entry, onSelect}: CardProps) {
    const Preview = entry.preview;

    return (
        <li className='DemoCompassShowcase__card'>
            <div
                inert={true}
                aria-hidden={true}
            >
                <ScaledPreview className={entry.previewSurface === 'sidebar' ? 'DemoCompassShowcase__cardPreview--sidebar' : undefined}>
                    <Preview/>
                </ScaledPreview>
            </div>
            <button
                type='button'
                className='DemoCompassShowcase__cardName'
                onClick={(event) => onSelect(entry, event.currentTarget)}
            >
                {entry.name}
            </button>
        </li>
    );
}
