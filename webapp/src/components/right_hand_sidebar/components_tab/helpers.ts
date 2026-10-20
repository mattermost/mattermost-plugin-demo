import type {CatalogEntry} from './types';

export function filterCatalog(entries: CatalogEntry[], query: string): CatalogEntry[] {
    const q = query.trim().toLowerCase();
    if (!q) {
        return entries;
    }
    return entries.filter((entry) => (
        entry.name.toLowerCase().includes(q) ||
        entry.category.toLowerCase().includes(q) ||
        entry.description.toLowerCase().includes(q)
    ));
}

export function importSnippet(entry: Pick<CatalogEntry, 'id'>): string {
    const exportName = entry.id.split('-').map((part) => part.charAt(0).toUpperCase() + part.slice(1)).join('');
    return `import {${exportName}} from '@mattermost/compass-ui/components/${entry.id}';`;
}
