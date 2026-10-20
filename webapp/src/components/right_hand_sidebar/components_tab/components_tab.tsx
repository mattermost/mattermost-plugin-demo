import React, {useEffect, useRef, useState} from 'react';

import Detail from './detail';
import Gallery from './gallery';
import type {CatalogEntry} from './types';

import './components_tab.css';

export default function ComponentsTab() {
    const [selected, setSelected] = useState<CatalogEntry | null>(null);
    const triggerRef = useRef<HTMLButtonElement | null>(null);

    useEffect(() => {
        if (!selected) {
            triggerRef.current?.focus();
        }
    }, [selected]);

    return (
        <div className='DemoCompassShowcase__nav'>
            <div
                className='DemoCompassShowcase__pane'
                inert={selected !== null}
                aria-hidden={selected ? true : undefined}
            >
                <Gallery
                    onSelect={(entry, trigger) => {
                        triggerRef.current = trigger;
                        setSelected(entry);
                    }}
                />
            </div>
            {selected && (
                <div className='DemoCompassShowcase__pane DemoCompassShowcase__pane--push'>
                    <Detail
                        key={selected.id}
                        entry={selected}
                        onBack={() => setSelected(null)}
                    />
                </div>
            )}
        </div>
    );
}
