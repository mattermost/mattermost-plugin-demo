import {matchesQuery, sections} from './surfaces';

const item = {
    name: 'Root component',
    api: 'registerRootComponent',
    where: 'Overlay shown above everything.',
};

describe('matchesQuery', () => {
    it('matches everything for an empty or whitespace query', () => {
        expect(matchesQuery(item, '')).toBe(true);
        expect(matchesQuery(item, '   ')).toBe(true);
    });

    it('matches name, API and description case-insensitively', () => {
        expect(matchesQuery(item, 'ROOT')).toBe(true);
        expect(matchesQuery(item, 'registerroot')).toBe(true);
        expect(matchesQuery(item, 'overlay')).toBe(true);
    });

    it('trims the query', () => {
        expect(matchesQuery(item, '  root  ')).toBe(true);
    });

    it('rejects unrelated queries', () => {
        expect(matchesQuery(item, 'toast')).toBe(false);
    });
});

describe('sections', () => {
    it('uses unique item names within each section', () => {
        for (const section of sections) {
            const names = new Set();
            for (const {name} of section.items) {
                expect(names.has(name)).toBe(false);
                names.add(name);
            }
        }
    });
});
