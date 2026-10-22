import React, {useEffect, useRef, useState} from 'react';

import Detail from './detail';
import Gallery from './gallery';
import type {CatalogEntry} from './types';

import './components_tab.css';

export default function ComponentsTab() {
    const [selected, setSelected] = useState<CatalogEntry | null>(null);
    const [exiting, setExiting] = useState(false);
    const triggerRef = useRef<HTMLButtonElement | null>(null);

    // The detail pane stays mounted while it slides out; it is removed on animationend.
    const handleBack = () => {
        if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
            setSelected(null);
            return;
        }
        setExiting(true);
    };

    useEffect(() => {
        if (exiting || (!selected && triggerRef.current)) {
            triggerRef.current?.focus({preventScroll: true});
        }
    }, [selected, exiting]);

    return (
        <div className='DemoCompassShowcase__nav'>
            <div
                className={`DemoCompassShowcase__pane DemoCompassShowcase__pane--base${selected && !exiting ? ' DemoCompassShowcase__pane--pushed' : ''}`}
                inert={selected !== null && !exiting}
                aria-hidden={selected && !exiting ? true : undefined}
            >
                <Gallery
                    onSelect={(entry, trigger) => {
                        triggerRef.current = trigger;
                        setExiting(false);
                        setSelected(entry);
                    }}
                />
            </div>
            {selected && (
                <div
                    className={`DemoCompassShowcase__pane DemoCompassShowcase__pane--push${exiting ? ' DemoCompassShowcase__pane--exiting' : ''}`}
                    inert={exiting}
                    onAnimationEnd={(event) => {
                        if (exiting && event.target === event.currentTarget) {
                            setSelected(null);
                            setExiting(false);
                        }
                    }}
                >
                    <Detail
                        key={selected.id}
                        entry={selected}
                        onBack={handleBack}
                    />
                </div>
            )}
        </div>
    );
}
